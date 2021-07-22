package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	log "github.com/go-pkgz/lgr"
	"github.com/the-NZA/acg-nikolaev/internal/app/acg"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config/acg_dev.json", "Path to config file")
}

func main() {
	flag.Parse()

	config := acg.NewConfig()
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	defer jsonFile.Close()

	configBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	s := acg.NewServer(config)
	if err := s.Start(); err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

}
