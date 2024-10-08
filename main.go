package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/safciplak/trustwallet/api"
	"github.com/safciplak/trustwallet/config"
	"github.com/safciplak/trustwallet/parser"
	"github.com/safciplak/trustwallet/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	var store storage.Storage

	switch cfg.StorageType {
	case "memory":
		store = storage.NewMemoryStorage()
	default:
		fmt.Println("Invalid storage type:", cfg.StorageType)
		return
	}

	rpcURL := "https://ethereum-rpc.publicnode.com"

	p := parser.NewParser(rpcURL, store)

	go func() {
		err := p.StartParsing()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}()

	apiServer := api.NewAPI(p)
	go apiServer.Start(":8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down gracefully...")
}
