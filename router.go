package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

func InitRouter() {
	r := mux.NewRouter()

	//GET
	r.HandleFunc("/get-player", GetPlayerRecord).Methods("GET")
	r.HandleFunc("/get-all-players", GetAllPlayersRecord).Methods("GET")

	//POST
	r.HandleFunc("/add-player", AddPlayer).Methods("POST")
	r.HandleFunc("/update-player-elo", UpdatePlayerElo).Methods("POST")
	r.HandleFunc("/delete-player", DeletePlayer).Methods("POST")
	r.HandleFunc("/calculate-match-result", CalculateMatchResult).Methods("POST")

	srv := &http.Server{
		Addr:    ":3333",
		Handler: r,
	}
	srv.ListenAndServe()
}

func GetAllPlayersRecord(w http.ResponseWriter, r *http.Request) {
	resp := []Record{}
	for _, record := range RecordMap {
		resp = append(resp, record)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := jsoniter.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func GetPlayerRecord(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Query().Get("player")

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := jsoniter.Marshal(RecordMap[player])
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record.Elo = DefaultElo
	RecordMap[record.Player] = record

	RewriteFileRecord()

	jsonResp, err := jsoniter.Marshal(record)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func UpdatePlayerElo(w http.ResponseWriter, r *http.Request) {
	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RecordMap[record.Player] = record

	RewriteFileRecord()

	jsonResp, err := jsoniter.Marshal(record)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//remove from global map first
	delete(RecordMap, record.Player)

	//then rewrite the file
	RewriteFileRecord()

	jsonResp, err := jsoniter.Marshal(record)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func CalculateMatchResult(w http.ResponseWriter, r *http.Request) {
	var matchResult MatchResult

	err := json.NewDecoder(r.Body).Decode(&matchResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	calculationResult := CalculateElo(matchResult)

	//update record map with new data
	for _, record := range calculationResult.Records {
		RecordMap[record.Player] = record
	}

	//update db
	RewriteFileRecord()

	jsonResp, err := jsoniter.Marshal(calculationResult)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
