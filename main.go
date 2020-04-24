package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/resources"
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

	conf, err := configuration.NewConfiguration(*configFile)
	if err != nil {
		log.Fatalf("read configuration: %s", err)
	}

	if *printConfig {
		conf.PrintConfiguration()
		return
	}

	storage, err := resources.NewStorage(conf.Storage, time.Duration(conf.Storage.Timeout)*time.Second)
	if err != nil {
		log.Fatalf("create storage: %s", err)
	}

	if *checkConfig {
		fmt.Println("Configuration checked")
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	server, err := NewServer(conf, storage)
	if err != nil {
		log.Fatalf("create server: %s", err)
	}

	server.Run(signalChan)
}
