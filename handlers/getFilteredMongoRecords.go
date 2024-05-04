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
		mrecs.getFilteredMongoRecords(w, r)
		return
	}
	// handle all else
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (mrecs *MongoRecords) getFilteredMongoRecords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle /POST request", r.URL)
	apiResponse := &data.ApiResponse{}

	// capture filter request from body
	filters := &data.FilterRequest{}
	err := filters.FromJSON(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		err = fmt.Errorf("unable to decode request payload from request body , err: %s", err)
		apiResponse.Code = http.StatusBadRequest
		apiResponse.Message = err.Error()
		apiResponse.ToJSON(w)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := filters.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get list of filtered mongo records
	mr, err := data.GetFilteredMongoRecords(database.DBconfig.ClientInstance, filters)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		err = fmt.Errorf("unable to fetch records from Mongodb, err: %s", err)
		apiResponse.Code = http.StatusInternalServerError
		apiResponse.Message = err.Error()
		apiResponse.ToJSON(w)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// package returned records into a response
	apiResponse.Code = 0
	apiResponse.Message = "Success"
	apiResponse.Records = mr
	w.Header().Set("Content-Type", "application/json")
	apiResponse.ToJSON(w)
}
