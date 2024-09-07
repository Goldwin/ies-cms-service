package mongo

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
	"go.mongodb.org/mongo-driver/mongo"
)

type searchHouseholdImpl struct {
	ctx context.Context
	db  *mongo.Database
}

// Execute implements queries.SearchHousehold.
func (s *searchHouseholdImpl) Execute(filter SearchHouseholdFilter) (SearchHouseholdResult, queries.QueryErrorDetail) {
	panic("unimplemented")
}

func NewSearchHousehold(ctx context.Context, db *mongo.Database) SearchHousehold {
	return &searchHouseholdImpl{
		ctx: ctx,
		db:  db,
	}
}
