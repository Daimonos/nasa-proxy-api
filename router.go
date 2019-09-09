package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/neos/{date}", GetNeosByDate)
	r.HandleFunc("/api/v1/marsweather", GetMarsWeather)
	r.HandleFunc("/api/v1/apod", GetApod)
	return r
}
