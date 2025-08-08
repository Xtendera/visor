package client

import (
	"ayode.org/visor/config"
	"ayode.org/visor/util"
	"ayode.org/visor/validations"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Client struct {
	cfg config.Config
	jar *cookiejar.Jar
	u   *url.URL
}

func New(cfg config.Config) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the cookiejar")
	}

	u, err := url.Parse(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the root URL")
	}

	c := Client{
		cfg,
		jar,
		u,
	}

	c.SetCookies(c.cfg.Jar)
	return &c, nil
}

func (c *Client) sendRequest(req *http.Request, logger *slog.Logger) (*http.Response, io.Reader, error) {
	client := &http.Client{
		Jar: c.jar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to close response stream: %s", err.Error()))
		}
	}(resp.Body)

	// Same as io.ReadAll(), but has less allocations therefore better performance
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
	return fmt.Errorf("invalid HTTP status recieved: %d", statusCode)
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
			return fmt.Errorf("could not marshal request body: %w", err)
		}

		*buffer = bytes.NewBuffer(jsonBody)
	}
	return nil
}

// prepareRequest creates and configures an HTTP request for the given endpoint
func (c *Client) prepareRequest(endpoint config.Endpoint) (*http.Request, error) {
	var reqBuff io.Reader
	err := marshalReqBody(endpoint.Body, &reqBuff)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	reqObj, err := http.NewRequest(endpoint.Method, c.cfg.Root+endpoint.Path, reqBuff)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequestHeaders(reqObj, endpoint)
	c.setRequestCookies(reqObj, endpoint)

	return reqObj, nil
}

// setRequestHeaders sets content-type and custom headers for the request
func (c *Client) setRequestHeaders(req *http.Request, endpoint config.Endpoint) {
	// Set content type based on body type
	if isJSONBody(endpoint.Body) {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "text/plain")
	}

	// Add custom headers (config-level)
	for _, header := range c.cfg.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	// Add custom headers (request-level)
	for _, header := range endpoint.Headers {
		req.Header.Set(header.Key, header.Value)
	}
}

// setRequestCookies sets cookies for the request
func (c *Client) setRequestCookies(req *http.Request, endpoint config.Endpoint) {
	// Add cookies from jar (config-level)
	c.SetCookies(endpoint.Jar)

	// Add cookies (request-level)
	c.SetReqCookies(req, endpoint.Cookies)
}

// processResponse handles the response processing and validation
func (c *Client) processResponse(resp *http.Response, respBody io.Reader, endpoint config.Endpoint) error {
	respBytes, err := io.ReadAll(respBody)
	if err != nil {
		return err
	}

	// Export to filesystem
	err = c.exportResponse(respBytes, endpoint)
	if err != nil {
		return err
	}

	// Validate status code
	err = validateStatus(&endpoint, uint16(resp.StatusCode))
	if err != nil {
		return err
	}

	// Validate schema if specified
	if endpoint.Schema != "" {
		respReader := bytes.NewReader(respBytes)
		err = validations.ValidateSchemaFromPath(respReader, endpoint.Schema)
		if err != nil {
			return fmt.Errorf("failed to validate response body: %w", err)
		}
	}

	return nil
}

func (c *Client) exportResponse(respBytes []byte, endpoint config.Endpoint) error {
	if c.cfg.SaveResponseDir != "" && endpoint.SaveResponse == "" {
		err := os.MkdirAll(c.cfg.SaveResponseDir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(c.cfg.SaveResponseDir, util.SanitizeFileName(endpoint.Name)+".json"))
		if err != nil {
			return err
		}
		_, err = f.WriteString(string(respBytes))
		if err != nil {
			err := f.Close()
			if err != nil {
				return err
			}
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	if endpoint.SaveResponse != "" {
		newPath := filepath.Dir(endpoint.SaveResponse)
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(endpoint.SaveResponse)
		if err != nil {
			return err
		}
		_, err = f.WriteString(string(respBytes))
		if err != nil {
			err := f.Close()
			if err != nil {
				return err
			}
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// executeEndpoint executes a single endpoint request
func (c *Client) executeEndpoint(endpoint config.Endpoint) {
	logger := slog.With("taskName", endpoint.Name, "path", endpoint.Path, "method", endpoint.Method)

	// Prepare the request
	reqObj, err := c.prepareRequest(endpoint)
	if err != nil {
		logger.Error("Failed to prepare request: " + err.Error())
		return
	}

	// Send the request
	start := time.Now()
	resp, respBody, err := c.sendRequest(reqObj, logger)
	elapsed := time.Since(start)

	if err != nil {
		logger.Error("Failed to send request: " + err.Error())
		return
	}

	// Process the response
	err = c.processResponse(resp, respBody, endpoint)
	if err != nil {
		logger.Error("Failed to process response: " + err.Error())
		return
	}

	logger.With("elapsed", elapsed).Info("Task Succeeded")
}

func (c *Client) Execute() {
	for _, endpoint := range c.cfg.Endpoints {
		c.executeEndpoint(endpoint)
	}
}
