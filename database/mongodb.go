package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDBClient(uri string, name string) error {
	clientOption := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//connect to mongoDB
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return err
	}
	//check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	log.Println("connected to monogodb!")
	//get the database
	MongoDB = client.Database(name)
	return nil
}
