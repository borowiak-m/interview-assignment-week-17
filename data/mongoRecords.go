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
	Key        string    `bson:"key" json:"key"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	Count      []int     `bson:"count" json:"-"`
	TotalCount int       `json:"totalCount"` // only for response payload
}

type Records []*Record

type ApiResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"msg"`
	Records []*Record `json:"records"`
}

type FilterRequest struct {
	StartDate JSONTime `json:"startDate"`
	EndDate   JSONTime `json:"endDate"`
	MinCount  int      `json:"minCount"`
	MaxCount  int      `json:"maxCount"`
}

type JSONTime struct {
	time.Time
}

const dateFormat = "2006-01-02"

// UnmarshalJSON parses JSON data into JSONTime.
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	t, err := time.Parse(`"`+dateFormat+`"`, s)
	if err != nil {
		return fmt.Errorf("date must be in format %s", dateFormat)
	}
	jt.Time = t
	return nil
}

// encoding a api response
func (resp *ApiResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(resp)
}

// decoding a single object
func (fr *FilterRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fr)
}

func (fr *FilterRequest) Validate() error {
	if fr.StartDate.Time.IsZero() {
		return fmt.Errorf("startDate must be provided and not be zero")
	}
	if fr.EndDate.Time.IsZero() {
		return fmt.Errorf("endDate must be provided and not be zero")
	}
	if fr.StartDate.Time.After(fr.EndDate.Time) {
		return fmt.Errorf("startDate must be before endDate")
	}
	if fr.MinCount < 0 {
		return fmt.Errorf("minCount must be non-negative")
	}
	if fr.MaxCount < 0 {
		return fmt.Errorf("maxCount must be non-negative")
	}
	if fr.MinCount > fr.MaxCount {
		return fmt.Errorf("minCount should be less than or equal to maxCount")
	}
	return nil
}

func GetFilteredMongoRecords(client *mongo.Client, filters *FilterRequest) (Records, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("maindatabase").Collection("records")

	// MongoDB query to filter by date range
	filter := bson.M{
		"createdAt": bson.M{
			"$gte": filters.StartDate.Time,
			"$lte": filters.EndDate.Time,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error running collection.Find(ctx, bson.M{}): %s", err)
	}
	defer cursor.Close(ctx)

	var records Records
	for cursor.Next(ctx) {
		var rec Record
		if err := cursor.Decode(&rec); err != nil {
			return nil, fmt.Errorf("error decoding record: %s", err)
		}

		totalCount := 0
		// sum all count entries between the min and max count from filters
		for _, count := range rec.Count {
			if count >= filters.MinCount && count <= filters.MaxCount {
				totalCount += count
			}
		}
		// add total count to obj
		rec.TotalCount = totalCount
		records = append(records, &rec)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error during cursor operation: %s", err)
	}

	return records, nil
}
