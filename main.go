package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

var cache *Cache

func main() {
	log.Println("Starting NASA API Proxy")
	cache = &Cache{}
	cache.InitializeClient("127.0.0.1:6379", 10, 30)
	log.Println(cache.Client.Ping())
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
