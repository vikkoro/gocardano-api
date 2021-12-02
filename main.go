package main

import (
	"github.com/vikkoro/gocardano-api/pkg/cardano"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"github.com/vikkoro/gocardano-api/pkg/files"
	"github.com/vikkoro/gocardano-api/pkg/handlers"
	"github.com/vikkoro/gocardano-api/pkg/parser"
	"github.com/vikkoro/gocardano-api/pkg/wallet"
)

func main() {

	// Read config file
	cfg := config.NewConfig("conf.json", ".env")

	// Create Services with DI
	cs := cardano.NewService(cfg)
	ws := wallet.NewService(cfg, cs)
	ps := parser.NewService(cfg)
	fs := files.NewService(cfg)

	handlers.NewRestService(cfg, ws, ps, fs)
}
