package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/daimonos/nasa"
	"github.com/daimonos/nasa/models"
	"github.com/gorilla/mux"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func GetNeosByDate(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	if date == "" {
		WriteError(errors.New("Missing Date"), w, http.StatusBadRequest)
		return
	}
	dParse, dParseErr := time.Parse("2006-01-02", date)
	sDate := dParse.Format("2006-01-02")
	if dParseErr != nil {
		WriteError(dParseErr, w, http.StatusBadRequest)
		return
	}
	bytes, err := DB.Get("NEOS", sDate)
	if err != nil {
		WriteError(err, w, http.StatusInternalServerError)
		return
	}
	if len(bytes) > 0 {
		var neos []models.Neo
		jsonParseErr := json.Unmarshal(bytes, &neos)
		if jsonParseErr != nil {
			WriteError(jsonParseErr, w, http.StatusInternalServerError)
			return
		}
		WriteJSON(neos, w, http.StatusOK)
		return
	}
	neo, err := nasa.GetFeed(dParse, dParse)
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
		return
	}
	DB.Set("NEOS", sDate, neo.NearEarthObjects[sDate])
	WriteJSON(neo.NearEarthObjects[sDate], w, http.StatusOK)
}

func GetApod(w http.ResponseWriter, r *http.Request) {
	apod, err := nasa.GetApod(time.Now(), true)
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
		return
	}
	WriteJSON(apod, w, http.StatusOK)
}

func GetMarsWeather(w http.ResponseWriter, r *http.Request) {
	weather, err := nasa.GetMarsWeather()
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
	}
	WriteJSON(weather, w, http.StatusOK)
}

func WriteError(err error, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(ErrorMessage{Message: err.Error()})
	if err != nil {
		WriteError(err, w, http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	return
}

func WriteJSON(payload interface{}, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(payload)
	if err != nil {
		WriteError(err, w, http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	return
}
