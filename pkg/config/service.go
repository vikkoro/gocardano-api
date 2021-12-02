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
	cfg := Configuration{}

	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	cfg.Passphrase = os.Getenv("PASSPHRASE")
	cfg.WalletId = os.Getenv("WALLET_ID")

	return &cfg
}
