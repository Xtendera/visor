package config

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"os"
)

type Config struct {
	Root      string     `json:"root"`
	Endpoints []Endpoint `json:"endpoints" validate:"min=1,dive"`
}

type Endpoint struct {
	Name         string      `json:"name" validate:"required"`
	Path         string      `json:"path" validate:"required"`
	Method       string      `json:"method" validate:"oneof=GET POST PUT HEAD DELETE OPTIONS PATCH"`
	Body         interface{} `json:"body"`
	AcceptStatus []uint16    `json:"acceptStatus" validate:"required"`
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

	// Validate using the validator package
	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		slog.Error("Could not validate configuration: " + err.Error())
		os.Exit(5)
	}

	return cfg
}
