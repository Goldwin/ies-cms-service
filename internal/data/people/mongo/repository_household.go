package mongo

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type personSlice []*entities.Person

type householdRepositoryImpl struct {
	ctx                 context.Context
	db                  *mongo.Database
	householdCollection *mongo.Collection
	personCollection    *mongo.Collection
}

// Delete implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) Delete(e *entities.Household) error {
	_, err := h.householdCollection.DeleteOne(h.ctx, bson.M{"_id": e.ID})
	return err
}

// Get implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) Get(id string) (*entities.Household, error) {
	var model HouseholdModel
	err := h.householdCollection.FindOne(h.ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return nil, err
	}
	household := toEntity(model)
	personList, err := h.listPersonByHouseholdID(id)

	if err != nil {
		return nil, err
	}
	household.Members = lo.Filter(personList, func(e *entities.Person, _ int) bool {
		return e.ID != household.HouseholdHead.ID
	})

	household.HouseholdHead = lo.FindOrElse(personList, nil, func(e *entities.Person) bool {
		return e.ID == household.HouseholdHead.ID
	})

	return household, nil
}

// List implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) List(householdIDList []string) ([]*entities.Household, error) {
	var householdList []HouseholdModel
	cursor, err := h.householdCollection.Find(h.ctx, bson.M{"_id": bson.M{"$in": householdIDList}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(h.ctx, &householdList); err != nil {
		return nil, err
	}
	households := lo.Map(householdList, func(model HouseholdModel, _ int) *entities.Household {
		result := toEntity(model)
		result.HouseholdHead = &entities.Person{ID: model.HouseholdHeadID}
		return result
	})

	personMaps, err := h.mapPersonByHouseholdIDList(householdIDList...)

	if err != nil {
		return nil, err
	}

	for _, household := range households {
		household.Members = lo.Filter(personMaps[household.ID], func(e *entities.Person, _ int) bool {
			return e.ID != household.HouseholdHead.ID
		})
		household.HouseholdHead = lo.FindOrElse(personMaps[household.ID], nil, func(item *entities.Person) bool {
			return item.ID == household.HouseholdHead.ID
		})
	}

	return households, nil
}

// Save implements repositories.HouseholdRepository.
func (h *householdRepositoryImpl) Save(e *entities.Household) (*entities.Household, error) {
	model := toHouseholdModel(e)

	_, err := h.householdCollection.UpdateOne(h.ctx, bson.M{"_id": e.ID}, bson.M{"$set": model}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	personIDList := lo.Map(e.Members, func(e *entities.Person, _ int) string { return e.ID })
	personIDList = append(personIDList, e.HouseholdHead.ID)

	_, err = h.personCollection.UpdateMany(h.ctx, bson.M{"_id": bson.M{"$in": personIDList}}, bson.M{"$set": bson.M{"householdId": e.ID}})

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (h *householdRepositoryImpl) listPersonByHouseholdID(householdID string) ([]*entities.Person, error) {
	var personList []PersonModel
	cursor, err := h.personCollection.Find(h.ctx, bson.M{"householdId": householdID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(h.ctx, &personList); err != nil {
		return nil, err
	}
	return lo.Map(personList, func(model PersonModel, _ int) *entities.Person {
		return model.toEntity()
	}), nil
}

func (h *householdRepositoryImpl) mapPersonByHouseholdIDList(householdIDList ...string) (map[string][]*entities.Person, error) {
	var personList []PersonModel
	cursor, err := h.personCollection.Find(h.ctx, bson.M{"householdId": householdIDList})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(h.ctx, &personList); err != nil {
		return nil, err
	}

	personModelByHouseholdID := lo.GroupBy(personList, func(model PersonModel) string {
		return model.HouseholdID
	})

	return lo.MapValues(personModelByHouseholdID, func(personList []PersonModel, _ string) []*entities.Person {
		return lo.Map(personList, func(model PersonModel, _ int) *entities.Person {
			return model.toEntity()
		})
	}), nil
}

func NewHouseholdRepository(ctx context.Context, db *mongo.Database) repositories.HouseholdRepository {
	return &householdRepositoryImpl{
		ctx:                 ctx,
		db:                  db,
		householdCollection: db.Collection(householdCollectionName),
		personCollection:    db.Collection(personCollectionName),
	}
}
