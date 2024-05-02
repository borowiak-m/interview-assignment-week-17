package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

func Connect(uri string) {
	fmt.Println("starting Connect...")
	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// instantiate client
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	// defer disconnecting client after instantiating
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	fmt.Println(" - - 1 no error after connecting...")
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Cannot connect to ping", err)
	}
	fmt.Println(" - - 1 no error after pinging...")
	fmt.Println("Successfully connected to MongoDB!")
}
