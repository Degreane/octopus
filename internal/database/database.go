package database

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDataFromCollection(uri string, dbName string, collectionName string, filter interface{}) (string, error) {
	// Connect to MongoDB
	client, err := ConnectDB(uri)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())

	// Get collection
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	// Get all results
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return "", err
	}

	// Convert to JSON string
	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func SetDataToCollection(uri string, dbName string, collectionName string, filter interface{}, data interface{}) (string, error) {
	// Connect to MongoDB
	client, err := ConnectDB(uri)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())

	// Get collection
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update or insert data
	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data}, opts)
	if err != nil {
		return "", err
	}

	// Convert result to JSON string
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func DelDataFromCollection(uri string, dbName string, collectionName string, filter interface{}) (string, error) {
	// Connect to MongoDB
	client, err := ConnectDB(uri)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())

	// Get collection
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete data
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return "", err
	}

	// Convert result to JSON string
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func InsertDataToCollection(uri string, dbName string, collectionName string, data interface{}) (string, error) {
	// Connect to MongoDB
	client, err := ConnectDB(uri)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())

	// Get collection
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert data
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return "", err
	}

	// Convert result to JSON string
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
