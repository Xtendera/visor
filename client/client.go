package client

import (
	"ayode.org/visor/config"
	"bytes"
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

func sendRequest(req *http.Request, responseBody *string) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if responseBody != nil {
		// Comparable to io.ReadAll(), but this has less allocations and therefore better performance.
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		*responseBody = buf.String()
	}

	return resp, nil
}

func (c Client) Execute() {
	for _, endpoint := range c.cfg.Endpoints {
		logger := slog.With("taskName", endpoint.Name)
		reqObj, err := http.NewRequest(endpoint.Method, c.cfg.Root+endpoint.Path, nil)
		if err != nil {
			logger.Error("Error when creating request: " + err.Error())
			continue
		}
		resp, err := sendRequest(reqObj, nil)

		if err != nil {
			logger.Error("Error when sending request: " + err.Error())
			continue
		}
		// Check if response code is successful
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			logger.Error("Unexpected status code: " + strconv.Itoa(resp.StatusCode))
		}

		logger.Info("Task Succeeded")
	}
}
