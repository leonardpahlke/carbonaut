package httpwrapper_test

import (
	"net/http"
	"testing"

	"carbonaut.dev/pkg/util/httpwrapper"
)

func TestSendHTTPRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		baseURL string
		wantErr bool
	}{
		{"Valid Endpoint POST", http.MethodPost, "https://httpbin.org/post", false},
		{"Valid Endpoint GET", http.MethodGet, "https://httpbin.org/get", false},
		{"Invalid Endpoint POST", http.MethodPost, "https://example.does-not-exist/", true},
		{"Invalid Endpoint GET", http.MethodGet, "https://example.does-not-exist/", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpwrapper.SendHTTPRequest(&httpwrapper.HTTPReqWrapper{
				Method:  tc.method,
				BaseURL: tc.baseURL,
			})

			if (err != nil) != tc.wantErr {
				t.Fatalf("SendHTTPRequest() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
				} else if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}
