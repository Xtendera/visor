package client

import (
	"ayode.org/visor/config"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type Client struct {
	cfg config.Config
}

func New(cfg config.Config) Client {
	c := Client{
		cfg,
	}
	return c
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c Client) Execute() {
	for _, endpoint := range c.cfg.Endpoints {
		resp, err := http.Get(c.cfg.Root + endpoint.Path)
		logger := slog.With("taskName", endpoint.Name)
		if err != nil {
			logger.Error("Error when sending request: " + err.Error())
			continue
		}

		// Check if response code is successful
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			logger.Error("Unexpected status code: " + strconv.Itoa(resp.StatusCode))
		}

		defer resp.Body.Close()
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Error when reading request: " + err.Error())
			continue
		}
		logger.Info("Task Succeeded")
	}
}
