package handlers

import (
	"fmt"
	"net/http"

	"github.com/borowiak-m/interview-assignment-week-17/data"
)

type InMemRecords struct {
}

func NewInMemRecords() *InMemRecords {
	return &InMemRecords{}
}

func (inmemrecs *InMemRecords) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle /GET
	if r.Method == http.MethodGet {
		inmemrecs.getInMemRecords(w, r)
		return
	}
	// handle /POST
	if r.Method == http.MethodPost {
		inmemrecs.addInMemRecord(w, r)
		return
	}
	// handle all else
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// /GET in memory records
func (inmemrecs *InMemRecords) getInMemRecords(w http.ResponseWriter, r *http.Request) {
	// log request details
	fmt.Println("Handle /GET request", r.URL)
	// get list of all records in memory
	lr := data.GetInMemRecords()
	w.Header().Set("Content-Type", "application/json")
	if err := lr.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshall json", http.StatusInternalServerError)
	}
}

// /POST in memory records
func (inmemrecs *InMemRecords) addInMemRecord(w http.ResponseWriter, r *http.Request) {
	// log request details
	fmt.Println("Handle /POST request", r.URL)
	// grab new item from request body
	newRec := &data.InMemRecord{}
	err := newRec.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to decode new record from request body", http.StatusBadRequest)
	}
	// log new record
	fmt.Printf("New record: %#v", newRec)
	// add new item to list
	data.AddRecordToMemory(newRec)
	// echo new object
	w.Header().Set("Content-Type", "application/json")
	if err = newRec.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshall json", http.StatusInternalServerError)
	}
}
