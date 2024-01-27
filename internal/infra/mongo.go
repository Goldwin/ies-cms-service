package infra

import (
	"context"
	"log"

	//"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	DbName   string `env:"MONGO_DB" yaml:"db_name" default:"church-management"`
	URI      string `env:"MONGO_URL" yaml:"hosts" default:"127.0.0.1:27017"`
	Username string `env:"MONGO_USERNAME" yaml:"username" default:""`
	Password string `env:"MONGO_PASSWORD" yaml:"password" default:""`
}

func NewMongoDatabase(m *MongoConfig) *mongo.Database {
	var mongoClient *mongo.Client
	if m == nil {
		return nil
	}

	log.Default().Printf("Initializing Mongo. Connecting to %s", m.URI)
	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(m.URI),
		options.Client().SetAuth(options.Credential{
			Username: m.Username,
			Password: m.Password,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to connect to mongo, caused by %v", err)
	}
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to mongo, caused by %v", err)
	}
	return mongoClient.Database(m.DbName)
}
