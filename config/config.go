package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Config struct {
	Root      string     `json:"root"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func Parse(cfgFile string) Config {
	raw, err := os.ReadFile(cfgFile)
	if err != nil {
		slog.Error("Could not read configuration: " + err.Error())
		os.Exit(3)
	}
	var cfg Config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		slog.Error("Could not unmarshal configuration: " + err.Error())
		os.Exit(4)
	}
	return cfg
}
