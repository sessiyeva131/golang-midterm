package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/store/{id}", returnRecord).Methods("GET")
	myRouter.HandleFunc("/store", createNewRecord).Methods("POST")
	myRouter.HandleFunc("/store/{id}", updateRecord).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type Store struct {
	mu     sync.Mutex
	stores map[string]string
}

type Record struct {
	Id    string `json:"Id"`
	Value string `json:"Value"`
}

var st = Store{
	stores: map[string]string{"1": "record1", "2": "record2", "3": "record3"},
}

var wg sync.WaitGroup

//UPDATING ELEMENTS IN MAP
var updateValue = func(k string, v string) {
	st.updateValueReference(k, v)
	wg.Done()
}

func (s *Store) updateValueReference(k string, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stores[k] = v
}

// RETURNS RECORD BY ID OR KEY
func returnRecord(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]

	for id, value := range st.stores {
		if id == key {
			var record Record
			record.Id = id
			record.Value = value
			json.NewEncoder(w).Encode(record)
		}
	}
}

//CREATES NEW RECORD
func createNewRecord(w http.ResponseWriter, request *http.Request) {
	var newRecord Record
	reqBody, _ := ioutil.ReadAll(request.Body)

	json.Unmarshal(reqBody, &newRecord)
	wg.Add(1)
	updateValue(newRecord.Id, newRecord.Value)
	wg.Wait()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRecord)
}

//UPDATES EXISTING RECORD
func updateRecord(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]

	var updatedRecord Record

	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &updatedRecord)

	for id := range st.stores {
		if id == key {
			wg.Add(1)
			updateValue(key, updatedRecord.Value)
			wg.Wait()

			var record Record
			record.Id = id
			record.Value = updatedRecord.Value
			json.NewEncoder(w).Encode(record)
		}
	}
}
