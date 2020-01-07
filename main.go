package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/storage"
	"github.com/danielkraic/knihomol/web"
	"github.com/danielkraic/knihomol/web/handlers"
	"github.com/spf13/pflag"
)

//TODO: http timeout

var (
	// Version will be set during build
	Version = ""
	// Commit will be set during build
	Commit = ""
	// Build will be set during build
	Build = ""
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	pflag.StringP("addr", "a", "0.0.0.0:80", "HTTP service address.")
	configFile := pflag.StringP("config", "c", "", "path to config file")
	checkConfig := pflag.Bool("config-check", false, "check configuration")
	printConfig := pflag.BoolP("print-config", "p", false, "print configuration")
	pflag.Parse()

	apiConfiguration, err := configuration.NewConfiguration(*configFile)
	if err != nil {
		log.Fatalf("failed to read configuration: %s", err)
	}

	if *printConfig {
		apiConfiguration.PrintConfiguration()
		return
	}

	apiStorage, err := storage.NewStorage(apiConfiguration.Storage, time.Duration(apiConfiguration.Storage.Timeout)*time.Second)
	if err != nil {
		log.Fatalf("failed to create storage: %s", err)
	}

	if *checkConfig {
		fmt.Println("Configuration checked")
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	web, err := web.NewWeb(&handlers.Version{Version: Version, Commit: Commit, Build: Build}, apiConfiguration, apiStorage)
	if err != nil {
		log.Fatalf("failed to create API: %s", err)
	}

	web.Run(signalChan)
}
