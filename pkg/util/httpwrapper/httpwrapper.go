/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"

	"github.com/google/go-querystring/query"
)

type HTTPReqWrapper struct {
	// HTTP method: example "GET" or "POST"
	Method string
	// Base url without any path: example "http://localhost:80"
	BaseURL string
	// Path of the http request which is added to the baseURL: example "/foo"
	Path string
	// Query parameters as struct which are set after the "Path" of the http request: example {foo: a, bar: b} -> "?foo=a&bar=b"
	QueryStruct interface{}
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
	slog.Info("prepare http request", "method", req.Method, "baseURL", req.BaseURL, "path", req.Path, "queryStruct", req.QueryStruct, "bodyStruct", req.BodyStruct, "headers", req.Headers)
	v, err := query.Values(req.QueryStruct)
	if err != nil {
		return nil, fmt.Errorf("unable to encode query params: %w", err)
	}
	url := fmt.Sprintf("%s%s?%s", req.BaseURL, req.Path, v.Encode())
	slog.Info("sending http request", "url", url)
	requestBodyBytes, err := json.Marshal(&req.BodyStruct)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %v: %w", req, err)
	}
	request, err := http.NewRequest(req.Method, url, bytes.NewBuffer(requestBodyBytes))
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	return &HTTPReqInfo{body, resp.StatusCode}, nil
}
