package infra

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type InfraConfig struct {
	RedisConfig RedisConfig `yaml:"redis"`
	MongoConfig MongoConfig `yaml:"mongo"`
}

type InfraComponent interface {
	Redis() redis.UniversalClient
	Mongo() *mongo.Database
}

type infraComponentImpl struct {
	redisClient redis.UniversalClient
	mongoClient *mongo.Database
}

func (i *infraComponentImpl) Redis() redis.UniversalClient {
	return i.redisClient
}

func (i *infraComponentImpl) Mongo() *mongo.Database {
	return i.mongoClient
}

func NewInfraComponent(config InfraConfig) InfraComponent {
	return &infraComponentImpl{
		redisClient: NewRedisClient(&config.RedisConfig),
		mongoClient: NewMongoDatabase(&config.MongoConfig),
	}
}
