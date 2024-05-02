// prom2json
// 	large parts of the logic copied from this project https://github.com/prometheus/prom2json
// 	with minor edits to turn it into a library and bake it into carbonaut

package promscraper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	dto "github.com/prometheus/client_model/go"
)

type Prom2Json struct {
	// URL of the endpoint that exposes metrics in prom format
	URL string
	// Client certificate file
	Cert string
	// Client certificate's key file
	Key string
	// Accept any certificate during TLS handshake. Insecure, use only for testing
	AcceptInvalidCert bool
}

// Collect scrapes the endpoint of the URL provided and transforms prom metrics into Prometheus Family struct types
func Collect(cfg Prom2Json) ([]*Family, error) {
	mfChan := make(chan *dto.MetricFamily, 1024)
	errChan := make(chan error, 1)

	transport, err := makeTransport(cfg.Cert, cfg.Key, cfg.AcceptInvalidCert)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	go func() {
		err := fetchMetricFamilies(cfg.URL, mfChan, transport)
		if err != nil {
			errChan <- err
		}
	}()

	result := []*Family{}
	for mf := range mfChan {
		result = append(result, newFamily(mf))
	}

	// Check if there are errors reported during fetchMetricFamilies
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return result, nil
}

// CollectJSON scrapes the endpoint of the URL provided and transforms prom metrics into JSON
func CollectJSON(cfg Prom2Json) ([]byte, error) {
	result, err := Collect(cfg)
	if err != nil {
		return nil, err
	}

	jsonText, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonText, nil
}

func makeTransport(
	certificate string, key string,
	skipServerCertCheck bool,
) (*http.Transport, error) {
	// Start with the DefaultTransport for sane defaults.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Conservatively disable HTTP keep-alives as this program will only
	// ever need a single HTTP request.
	transport.DisableKeepAlives = true
	// Timeout early if the server doesn't even return the headers.
	transport.ResponseHeaderTimeout = time.Minute
	tlsConfig := &tls.Config{InsecureSkipVerify: skipServerCertCheck}
	if certificate != "" && key != "" {
		cert, err := tls.LoadX509KeyPair(certificate, key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	transport.TLSClientConfig = tlsConfig
	return transport, nil
}
