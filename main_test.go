package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/borowiak-m/interview-assignment-week-17/database"
	"github.com/borowiak-m/interview-assignment-week-17/handlers"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env:        []string{"MONGO_INITDB_DATABASE=testdb"},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var db *mongo.Client
	if err := pool.Retry(func() error {
		var err error
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		db, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		return db.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// setup env
	os.Setenv("MONGO_URI", fmt.Sprintf("mongodb://localhost:%s/testdb", resource.GetPort("27017/tcp")))
	os.Setenv("DATABASE", "testdb")
	os.Setenv("COLLECTION", "records")

	// init db
	database.InitDatabaseConn(
		os.Getenv("MONGO_URI"),
		os.Getenv("DATABASE"),
		os.Getenv("COLLECTION"))

	// run tests
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestFetchMongoRecordsSuccess(t *testing.T) {
	// create request body
	reqBody := strings.NewReader(`{
		"startDate":"2024-01-01",
		"endDate":"2024-12-31",
		"minCount": 100,
		"maxCount": 800
	}`)

	req, err := http.NewRequest("POST", "/fetchMongoRecords", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// create a response recorder and handler
	rr := httptest.NewRecorder()
	handler := handlers.NewMongoRecords()

	// call handler
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code: got %v want %v", status, http.StatusOK)
	}

	// check response payload
	expected := `{"code":0, "msg":"Success", "records":[...]}`
	if !strings.Contains(rr.Body.String(), `"code":0`) || !strings.Contains(rr.Body.String(), `"msg":"Success"`) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFetchMongoRecordsBadRequest(t *testing.T) {
	// create request with missing fields
	reqBody := strings.NewReader(`{"startDate":"2024-01-01"}`)
	req, err := http.NewRequest("POST", "/fetchMongoRecords", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// create a response recorder and handler
	rr := httptest.NewRecorder()
	handler := handlers.NewMongoRecords()

	// call handler
	handler.ServeHTTP(rr, req)

	// verify status code
	expected := http.StatusBadRequest
	if status := rr.Code; status != expected {
		t.Errorf("handler returned incorrect status code: got %v want %v",
			status, expected)
	}
}

// testing /GET in memory records enpoint
func TestGetInMemRecords(t *testing.T) {
	req, err := http.NewRequest("GET", "/inmemory", nil)
	if err != nil {
		t.Fatal(err)
	}

	// test recorder and handler
	rr := httptest.NewRecorder()
	handler := handlers.NewInMemRecords()

	// call handler
	handler.ServeHTTP(rr, req)

	// if response not ok
	expectedStatus := http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	// check response payload

	var got []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &got)
	if err != nil {
		t.Fatal("could not unmarshall response")
	}

	want := []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		{Key: "active-tabs", Value: "getir"},
		{Key: "inactive-tabs", Value: "getout!"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

// test /POST in memory records endpoint
func TestPostInMemRecords(t *testing.T) {
	newRecord := `{"key":"new-key","value":"new-value"}`
	req, err := http.NewRequest("POST", "/inmemory", bytes.NewBufferString(newRecord))
	if err != nil {
		t.Fatal(err)
	}

	// new recorder and handler
	rr := httptest.NewRecorder()
	handler := handlers.NewInMemRecords()

	// call handler
	handler.ServeHTTP(rr, req)

	// if response not ok
	expectedStatus := http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	// check response payload
	expectedPayload := newRecord
	if !strings.Contains(rr.Body.String(), expectedPayload) {
		t.Errorf("handler returned unexpected boty: got %v want %v",
			rr.Body.String(), expectedPayload)
	}

	// check if record added to in memory persistence
	req, _ = http.NewRequest("GET", "/inmemory", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), expectedPayload) {
		t.Errorf("new record was not added correctly, body: %v", rr.Body.String())
	}

}

// test invalid request payload for /POST
func TestPostInMemRecordsInvalidJSON(t *testing.T) {
	badJSON := `{"key":123}` // key should be string
	req, err := http.NewRequest("POST", "/inmemory", bytes.NewBufferString(badJSON))
	if err != nil {
		t.Fatal(err)
	}

	// new recorder and handler
	rr := httptest.NewRecorder()
	handler := handlers.NewInMemRecords()

	// call handler
	handler.ServeHTTP(rr, req)

	// check response status
	expectedStatus := http.StatusBadRequest
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code for invalid request payload: got %v want %v",
			status, expectedStatus)
	}

}
