package main

import (
	"github.com/vikkoro/gocardano-api/pkg/cardano"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"github.com/vikkoro/gocardano-api/pkg/files"
	"github.com/vikkoro/gocardano-api/pkg/handlers"
	"github.com/vikkoro/gocardano-api/pkg/parser"
	"github.com/vikkoro/gocardano-api/pkg/wallet"
	"log"
	"os"
)

func main() {

	// Setup log file
	logging("info.log")

	// Read config file
	cfg := config.NewConfig("conf.json", ".env")

	// Create Services with DI
	cs := cardano.NewService(cfg)
	ws := wallet.NewService(cfg, cs)
	ps := parser.NewService(cfg)
	fs := files.NewService(cfg)

	handlers.NewRestService(cfg, ws, ps, fs)
}

// Setup log file
func logging(logFile string) {

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = file.Close()
	}()

	log.SetOutput(file)
	log.Print("Logging to a file in Go!")
}
