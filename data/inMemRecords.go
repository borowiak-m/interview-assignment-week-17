package data

import (
	"encoding/json"
	"io"
)

type InMemRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InMemRecords []*InMemRecord

// encoding a list of objects
func (inmemrecs *InMemRecords) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(inmemrecs)
}

// decoding a single object
func (inmemrec *InMemRecord) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(inmemrec)
}

// encoding a single of objects
func (inmemrec *InMemRecord) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(inmemrec)
}

var inMemRecords = InMemRecords{
	&InMemRecord{
		Key:   "active-tabs",
		Value: "getir",
	},
	&InMemRecord{
		Key:   "inactive-tabs",
		Value: "getout!",
	},
}

func GetInMemRecords() InMemRecords {
	return inMemRecords
}

func AddRecordToMemory(inMemRec *InMemRecord) {
	inMemRecords = append(inMemRecords, inMemRec)
}
