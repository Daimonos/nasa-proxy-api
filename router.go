package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/neos/{date}", GetNeosByDate).Methods("GET")
	r.HandleFunc("/api/v1/marsweather", GetMarsWeather).Methods("GET")
	r.HandleFunc("/api/v1/apod", GetApod).Methods("GET")
	r.HandleFunc("/api/v1/apod/{date}", GetApod).Methods("GET")
	return r
}
