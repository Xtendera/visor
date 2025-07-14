package client

import (
	"ayode.org/visor/config"
	"ayode.org/visor/validations"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"time"
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

func sendRequest(req *http.Request) (*http.Response, io.Reader, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, buf, nil
}

func validateStatus(endpoint *config.Endpoint, statusCode uint16) error {
	if endpoint.AcceptStatus == nil || len(endpoint.AcceptStatus) == 0 {
		return nil
	}
	for _, currStatus := range endpoint.AcceptStatus {
		if currStatus == statusCode {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid HTTP status recieved: %d", statusCode))
}

func isJSONBody(body interface{}) bool {
	switch body.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		// Check for struct types (since it technically doesn't have an explicit type in the Golang system)
		t := reflect.TypeOf(body)
		return t != nil && t.Kind() == reflect.Struct
	}
}

func marshalReqBody(body interface{}, buffer *io.Reader) error {
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not marshal request body: %s", err.Error()))
		}

		*buffer = bytes.NewBuffer(jsonBody)
	}
	return nil
}

func (c Client) Execute() {
	for _, endpoint := range c.cfg.Endpoints {
		logger := slog.With("taskName", endpoint.Name, "path", endpoint.Path, "method", endpoint.Method)

		var reqBuff io.Reader
		err := marshalReqBody(endpoint.Body, &reqBuff)
		if err != nil {
			logger.Error("Error when marshaling request body:" + err.Error())
			continue
		}
		reqObj, err := http.NewRequest(endpoint.Method, c.cfg.Root+endpoint.Path, reqBuff)

		if err != nil {
			logger.Error("Error when creating request: " + err.Error())
			continue
		}

		if isJSONBody(endpoint.Body) {
			reqObj.Header.Set("Content-Type", "application/json")
		} else {
			reqObj.Header.Set("Content-Type", "text/plain")
		}

		start := time.Now()
		resp, respBody, err := sendRequest(reqObj)
		elapsed := time.Now().Sub(start)

		if err != nil {
			logger.Error("Error when sending request: " + err.Error())
			continue
		}

		// Check if response code is the expected value
		err = validateStatus(&endpoint, uint16(resp.StatusCode))
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if endpoint.SchemaFile != "" {
			err = validations.ValidateSchemaFromFile(respBody, endpoint.SchemaFile)
			if err != nil {
				logger.Error("Failed to validate response body: " + err.Error())
				continue
			}

		}
		logger.With("elapsed", elapsed).Info("Task Succeeded")
	}
}
