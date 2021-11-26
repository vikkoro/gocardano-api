package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	Handlers "github.com/vikkoro/gocardano-api/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	configuration := GetConfig("conf.json")

	// Create a mux router
	r := mux.NewRouter()

	// We will define a single endpoint
	r.Handle("/api/v1/cardano/{module}", Handlers.ClientHandler{Configuration: configuration})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))

	// Listen to port 8080 for incoming REST calls
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetConfig(path string) Handlers.Configuration {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error: ", err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	configuration := Handlers.Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error: ", err)
	}

	return configuration
}
