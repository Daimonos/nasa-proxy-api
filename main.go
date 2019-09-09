package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	log.Println("Starting NASA API Proxy")
	DB.Open("nasa.db")
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
