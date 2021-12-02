package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// NewService constructor of the default service.
func NewConfig(conf string, env string) *Configuration {
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	file, err := os.Open(conf)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	configuration.Passphrase = os.Getenv("PASSPHRASE")

	return &configuration
}
