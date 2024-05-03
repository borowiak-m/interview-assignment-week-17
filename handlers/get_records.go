package handlers

import (
	"fmt"
	"net/http"

	"github.com/borowiak-m/interview-assignment-week-17/data"
	"github.com/borowiak-m/interview-assignment-week-17/database"
)

// /GET records from Mongo
// 1. request payload:
// ---JSON: 4 fields
// startDate date
// endDate date
// minCount int
// maxCount int

// 2. response payload:
// ---JSON: 3 fields
// code: status code (0 for success)
// msg: status description (success)
// records: array of filtered items including fields
// - key
// - createdAt
// - totalCount = sum of count in the document/table

type MongoRecords struct {
}

func NewMongoRecords() *MongoRecords {
	return &MongoRecords{}
}

func (mrecs *MongoRecords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle /POST
	if r.Method == http.MethodPost {
		mrecs.getMongoRecords(w, r)
		return
	}
	// handle all else
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (mrecs *MongoRecords) getMongoRecords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle /POST request", r.URL)
	// get list of all mongo records
	mr, err := data.GetMongoRecords(database.DBconfig.ClientInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to fetch records from Mongodb, err: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = mr.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshall json", http.StatusInternalServerError)
	}
}
