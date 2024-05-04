package handlers

import (
	"fmt"
	"net/http"

	"github.com/borowiak-m/interview-assignment-week-17/data"
	"github.com/borowiak-m/interview-assignment-week-17/database"
)

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

	// capture filter request from body
	filters := &data.FilterRequest{}
	if err := filters.FromJSON(r.Body); err != nil {
		respondBadRequest(w, fmt.Sprintf("unable to decode request payload from request body , err: %s", err))
		return
	}

	// validate request params
	if err := filters.Validate(); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	// get list of filtered mongo records
	mr, err := data.GetFilteredMongoRecords(database.DBconfig.ClientInstance, filters)
	if err != nil {
		respondInternalServerError(w, fmt.Sprintf("unable to fetch records from Mongodb, err: %s", err))
		return
	}

	// package returned records into a response
	respondSuccess(w, mr)
}
