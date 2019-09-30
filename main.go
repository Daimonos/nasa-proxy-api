package main

import (
	"io"
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
	handler := handlers.CORS()(r)
	logFile, logFileError := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if logFileError != nil {
		log.Fatalf(logFileError.Error())
	}
	httpLogFile, httpFileError := os.OpenFile("http_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if httpFileError != nil {
		log.Fatal(httpFileError.Error())
	}
	logWriter := io.MultiWriter(os.Stdout, logFile)
	httpLogWriter := io.MultiWriter(os.Stdout, httpLogFile)
	handler = handlers.ProxyHeaders(handler)
	handler = handlers.LoggingHandler(httpLogWriter, handler)
	log.SetOutput(logWriter)
	log.Fatal(http.ListenAndServe(getEnvVariable("NASA_PORT", ":80"), handler))
}

func getEnvVariable(name, defaultValue string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
