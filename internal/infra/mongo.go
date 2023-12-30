package infra

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Hosts    string `yaml:"hosts" default:"127.0.0.1:27017"`
	Username string `yaml:"username" default:""`
	Password string `yaml:"password" default:""`
}

func NewMongoClient(m *MongoConfig) *mongo.Client {
	var mongoClient *mongo.Client
	if m == nil {
		return nil
	}

	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().SetHosts(strings.Split(m.Hosts, ",")),
		options.Client().SetAuth(options.Credential{
			Username: m.Username,
			Password: m.Password,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to connect to mongo, caused by %v", err)
	}
	return mongoClient
}
