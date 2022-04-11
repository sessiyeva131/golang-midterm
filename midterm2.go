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
	Records = []Record{
		Record{Id: "1", Value: "Hello"},
		Record{Id: "2", Value: "Hello"},
	}
	handleRequests()
}

type Store struct {
	mu     sync.Mutex
	stores map[string]string
}

type Record struct {
	Id    string `json:"Id"`
	Value string `json:"Name"`
}

var s = Store{
	stores: map[string]string{"id": "record1", "id2": "record2"},
}

var wg sync.WaitGroup

var updateValue = func(key string, newValue string) {
	c.updateValueReference(key, newValue)
	wg.Done()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/store/{id}", returnRecord)
	myRouter.HandleFunc("/store", createNewRecord).Methods("POST")
	router.HandleFunc("/store/{id}", updateRecord).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for id, value := range s.stores {
		if id == key {
			var el Record
			el.Id = id
			el.Value = value
			json.NewEncoder(w).Encode(el)
		}
	}
}

func createNewRecord(w http.ResponseWriter, r *http.Request) {
	var newEvent Record
	reqBody, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(reqBody, &newEvent)
	wg.Add(1)
	updateValue(newEvent.Id, newEvent.Value)
	wg.Wait()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var updatedEvent Record

	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updatedEvent)

	for id := range s.stores {
		if id == key {
			wg.Add(1)
			updateValue(key, updatedEvent.Value)
			wg.Wait()

			var el Record
			el.Id = id
			el.Value = updatedEvent.Value
			json.NewEncoder(w).Encode(updatedElement)
		}
	}
}
