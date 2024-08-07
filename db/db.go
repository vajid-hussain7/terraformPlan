package db

import (
	"context"
	"fmt"
	"jsonTest/model"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Collection

func InitDB() error {
	connectionString := "mongodb://localhost:27017"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	DB = client.Database("sandbox").Collection("json")
	return nil
}

func InsertResourse(data model.ResourseChanges) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := DB.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	fmt.Println("data inserted  id is ", result.InsertedID)

	return nil
}
