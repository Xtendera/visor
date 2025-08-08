package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/url"
	"os"
)

type Config struct {
	Root            string     `json:"root" validate:"required"`
	Headers         []Header   `json:"headers"`
	Jar             []Cookie   `json:"jar"`
	Endpoints       []Endpoint `json:"endpoints" validate:"min=1,dive"`
	SaveResponseDir string     `json:"saveResponseDir"`
}

type Endpoint struct {
	Name         string      `json:"name" validate:"required"`
	Path         string      `json:"path" validate:"required"`
	Method       string      `json:"method" validate:"oneof=GET POST PUT HEAD DELETE OPTIONS PATCH"`
	Headers      []Header    `json:"headers"`
	Cookies      []Cookie    `json:"cookies"`
	Jar          []Cookie    `json:"jar"`
	Body         interface{} `json:"body"`
	AcceptStatus []uint16    `json:"acceptStatus" validate:"required"`
	SaveResponse string      `json:"saveResponse"`
	Schema       string      `json:"schema"`
}

type Header struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// Cookie Although the struct http.Cookie does exist, it does not contain JSON serialization keys, hence the need for a custom type
type Cookie struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// isValidAbsoluteURL checks if the given root URL is valid, and contains a scheme (e.g. https://), and is absolute (contains a host)
func isValidAbsoluteURL(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		return errors.New("invalid scheme")
	}
	if u.Host == "" {
		return errors.New("invalid absolute host")
	}
	return nil
}

func Parse(cfgFile string) Config {
	raw, err := os.ReadFile(cfgFile)
	if err != nil {
		slog.Error("Failed to read configuration: " + err.Error())
		os.Exit(3)
	}
	var cfg Config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		slog.Error("Failed to unmarshal configuration: " + err.Error())
		os.Exit(4)
	}

	// Validate using the validator package
	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		slog.Error("Failed to validate configuration: " + err.Error())
		os.Exit(5)
	}

	// Validate root URL
	err = isValidAbsoluteURL(cfg.Root)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to validate root URL \"%s\": %s", cfg.Root, err))
		os.Exit(6)
	}

	return cfg
}
