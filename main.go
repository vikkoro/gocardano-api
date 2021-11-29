package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	Handlers "github.com/vikkoro/gocardano-api/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	configuration, err := GetConfig("conf.json")
	if err != nil {
		log.Fatal("main:", err)
	}

	// Create a mux router
	r := mux.NewRouter()

	// We will define a single endpoint
	r.Handle("/api/v1/cardano/{module}", Handlers.ClientHandler{Configuration: configuration})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))

	// Listen to port 8080 for incoming REST calls
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetConfig(path string) (Handlers.Configuration, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return Handlers.Configuration{}, err
	}

	file, err := os.Open(path)

	if err != nil {
		return Handlers.Configuration{}, err
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	configuration := Handlers.Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		return Handlers.Configuration{}, err
	}

	configuration.Passphrase = os.Getenv("PASSPHRASE")

	return configuration, nil
}
