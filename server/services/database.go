package services

import (
	"context"
	"fmt"
	"github.com/mundanelizard/koyi/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Client
var collectionCache = map[string]*mongo.Collection{}

func init() {
	MongoClient = connectDatabase(config.MongoUri)
}

func connectDatabase(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	searchName := fmt.Sprintf("%s:%s", databaseName, collectionName)
	collection, ok := collectionCache[searchName]

	if !ok {
		collection = MongoClient.Database(databaseName).Collection(collectionName)
		collectionCache[searchName] = collection
	}

	return collection
}

func GetDatabase(databaseName string) *mongo.Database {
	return MongoClient.Database(databaseName)
}
