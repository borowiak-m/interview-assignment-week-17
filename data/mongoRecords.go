package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Record struct {
	Key       string    `bson:"key"`
	CreatedAt time.Time `bson:"createdAt"`
	StartDate time.Time `bson:"startDate"`
	EndDate   time.Time `bson:"endDate"`
	Count     []int     `bson:"count"`
}

func FetchRecords(client *mongo.Client) ([]Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("maindatabase").Collection("records")

	var records []Record
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil

}
