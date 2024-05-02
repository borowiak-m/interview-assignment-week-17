package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbConfig struct {
	ClientInstance *mongo.Client
	Database       string
	Collection     string
}

var DBconfig *dbConfig

func InitDatabaseConn(uri, database, collection string) {
	DBconfig = newDbConfig(uri, database, collection)
}

func newDbConfig(uri, database, collection string) *dbConfig {
	if DBconfig == nil {
		DBconfig = &dbConfig{
			Database:   database,
			Collection: collection,
		}
		connect(uri)
	}
	return DBconfig
}

func connect(uri string) {
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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Cannot connect to ping: ", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	DBconfig.ClientInstance = client
}
