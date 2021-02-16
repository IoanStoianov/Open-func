package repo

import (
	"context"
	"log"

	"github.com/IoanStoianov/Open-func/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepo implementation
type MongoRepo struct {
	mongo.Client
}

// CreateMongoClient - MongoRepo factory
func CreateMongoClient() (*MongoRepo, error) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return &MongoRepo{*client}, nil
}

// AddRecord does what it says
func (r *MongoRepo) AddRecord(result *types.FuncResult) error {
	collection := r.Database("open-func", &options.DatabaseOptions{}).Collection("results")

	res, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		return err
	}

	log.Printf("Record added with id %s", res.InsertedID)

	return nil
}

// GetRecords fetches last N records by name
func (r *MongoRepo) GetRecords(name string, count int64) ([]*types.FuncResult, error) {
	collection := r.Database("open-func", &options.DatabaseOptions{}).Collection("results")

	filter := bson.D{primitive.E{Key: "name", Value: name}}
	opts := options.FindOptions{
		Sort:  bson.D{primitive.E{Key: "_id", Value: -1}},
		Limit: &count,
	}

	cursor, err := collection.Find(context.Background(), filter, &opts)
	if err != nil {
		return nil, err
	}

	var results []*types.FuncResult

	for cursor.Next(context.Background()) {
		var result types.FuncResult
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(context.Background())

	return results, nil
}
