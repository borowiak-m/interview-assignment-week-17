package handlers

import "net/http"

type RecsInMem struct {
}

func NewRecsInMem() *RecsInMem {
	return &RecsInMem{}
}

func (inmem *RecsInMem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Response from RecsInMem"))
}

// /GET from map

// 3. GET request payload:
// key param in query param
// sample: localhost/in-memory?key=active-tabs
// 4. GET response payload
// ---JSON 2 fields
// key
// value
