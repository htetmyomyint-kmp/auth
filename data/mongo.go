package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "leave-tracker"
	collection = "users"
)

type MongoDBClient struct {
	Logger *log.Logger
	Client *mongo.Client
}

func NewMongoDBClient(l *log.Logger) *MongoDBClient {
	return &MongoDBClient{
		Logger: l,
	}
}

func (m *MongoDBClient) InitDB(ctx context.Context, config DBConfig) error {
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnString))
	if err != nil {
		return err
	}
	m.Client = cli
	return nil
}

func (m *MongoDBClient) GetUserByID(id string) (User, error) {
	return User{}, nil
}
func (m *MongoDBClient) CreateUser(user User) (User, error) {
	coll := m.Client.Database(database).Collection(collection)
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (m *MongoDBClient) GetUserByEmail(email string) (User, error) {
	var user User
	coll := m.Client.Database(database).Collection(collection)
	res := coll.FindOne(context.TODO(), map[string]interface{}{"email": email})

	if res.Err() != nil {
		return User{}, res.Err()
	}

	if err := res.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
