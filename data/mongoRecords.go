package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type Records []*Record

// encoding a list of objects
func (recs *Records) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(recs)
}

// decoding a single object
func (rec *Records) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(rec)
}

// encoding a single of objects
func (rec *Record) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rec)
}

func GetMongoRecords(client *mongo.Client) (Records, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("maindatabase").Collection("records")

	var records Records
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error running collection.Find(ctx, bson.M{}): %s", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			return nil, fmt.Errorf("error decoding record: %s", err)
		}
		records = append(records, &record)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during cursor operation: %s", err)
	}

	return records, nil

}
