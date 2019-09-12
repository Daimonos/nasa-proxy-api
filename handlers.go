package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/daimonos/nasa"
	"github.com/daimonos/nasa/models"
	"github.com/gorilla/mux"
)

const DAYSECONDS = 60 * 60 * 24
const HOURSECONDS = 60 * 60

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
	bytes, err := cache.Get("NEOS:" + sDate)
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
	log.Println("Getting new NEO's")
	neo, err := nasa.GetFeed(dParse, dParse)
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
		return
	}
	cache.Set("NEOS:"+sDate, neo.NearEarthObjects[sDate], DAYSECONDS)
	WriteJSON(neo.NearEarthObjects[sDate], w, http.StatusOK)
}

func GetApod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := time.Now()
	var err error
	if vars["date"] != "" {
		t, err = time.Parse("2006-01-02", vars["date"])
		if err != nil {
			WriteError(err, w, http.StatusInternalServerError)
			return
		}
	}
	key := t.Format("2006-01-02")
	bytes, err := cache.Get("APOD:" + key)
	if err != nil {
		WriteError(err, w, http.StatusInternalServerError)
		return
	}
	if len(bytes) > 0 {
		var apod *models.Apod
		jsonParseErr := json.Unmarshal(bytes, &apod)
		if jsonParseErr != nil {
			WriteError(jsonParseErr, w, http.StatusInternalServerError)
			return
		}
		WriteJSON(apod, w, http.StatusOK)
		return
	}
	apod, err := nasa.GetApod(t, true)
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
		return
	}
	cache.Set("APOD:"+key, apod, DAYSECONDS)
	WriteJSON(apod, w, http.StatusOK)
}

func GetMarsWeather(w http.ResponseWriter, r *http.Request) {
	bytes, err := cache.Get("WEATHER")
	if err != nil {
		WriteError(err, w, http.StatusInternalServerError)
		return
	}
	if len(bytes) > 0 {
		var weather models.MarsWeatherResp
		err := json.Unmarshal(bytes, &weather)
		if err != nil {
			WriteError(err, w, http.StatusInternalServerError)
			return
		}
		WriteJSON(weather, w, http.StatusOK)
		return
	}
	weather, err := nasa.GetMarsWeather()
	if err != nil {
		WriteError(err, w, http.StatusBadRequest)
	}
	cache.Set("WEATHER", weather, HOURSECONDS)
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
