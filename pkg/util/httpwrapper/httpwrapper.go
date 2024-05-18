package httpwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type HTTPReqWrapper struct {
	// HTTP method: example "GET" or "POST"
	Method string
	// Base url without any path: example "http://localhost:80"
	BaseURL string
	// Path of the http request which is added to the baseURL: example "/foo"
	Path string
	// Query parameters as struct which are set after the "Path" of the http request: example {foo: a, bar: b} -> "?foo=a&bar=b"
	Query string
	// Body as struct which gets send via POST requests: example {foo: a, bar: b} -> "{\"foo\": \"a\", \"bar\": \"b\"}"
	BodyStruct interface{}
	// list of http headers: example {Key: "Content-Type", Val: "application/json"}
	Headers map[string]string
}

type HTTPReqInfo struct {
	Body       []byte
	StatusCode int
}

// SendHTTPRequest is a wrapper function for sending http requests
// response: request body, status code, error
func SendHTTPRequest(req *HTTPReqWrapper) (*HTTPReqInfo, error) {
	// ? headerKey's are printed since headers are often used to transport secrets and accesskeys
	headerKey := []string{}
	for key := range req.Headers {
		headerKey = append(headerKey, key)
	}

	slog.Debug("prepare http request", "method", req.Method, "baseURL", req.BaseURL, "path", req.Path, "queryStruct", req.Query, "bodyStruct", req.BodyStruct, "header keys", headerKey)
	url := fmt.Sprintf("%s%s%s", req.BaseURL, req.Path, req.Query)
	slog.Info("sending http request", "url", url)

	var reader io.Reader
	if req.BodyStruct != nil {
		requestBodyBytes, err := json.Marshal(&req.BodyStruct)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal request body: %v: %w", req, err)
		}
		reader = bytes.NewBuffer(requestBodyBytes)
	}

	request, err := http.NewRequest(req.Method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("unable to create new %s request: %w", req.Method, err)
	}
	for k := range req.Headers {
		request.Header.Set(k, req.Headers[k])
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to perform %s request: %w", req.Method, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Warn("failed to close response body", "error", closeErr)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	return &HTTPReqInfo{body, resp.StatusCode}, nil
}
