package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

var cache *Cache

func main() {
	log.Println("Starting NASA API Proxy")
	key := getEnvVariable("NASA_API_KEY", "")
	if key == "" {
		log.Fatal("NASA API Key must be provided")
	}
	cache = &Cache{}
	cache.InitializeClient(getEnvVariable("NASA_REDIS_URL", "127.0.0.1:6379"), 10, 30)
	log.Println(cache.Client.Ping())
	r := NewRouter()
	log.Fatal(http.ListenAndServe(getEnvVariable("NASA_PORT", ":80"), handlers.CORS()(r)))
}

func getEnvVariable(name, defaultValue string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
