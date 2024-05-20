package dynres

import (
	"carbonaut.dev/pkg/provider/account/project/resource"
	"carbonaut.dev/pkg/provider/plugin"
)

type Config struct {
	Plugin plugin.Kind `json:"plugin" yaml:"plugin"`
	// Endpoint that is accessed to collec the data.
	// The IPv4 address will be collected from the static data thats looked up.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// Client certificate file
	Cert string `json:"cert" yaml:"cert"`
	// Client certificate's key file
	Key string `json:"key" yaml:"key"`
	// Accept any certificate during TLS handshake. Insecure, use only for testing
	AcceptInvalidCert bool `json:"accept_invalid_cert" yaml:"accept_invalid_cert"`
}

type Provider interface {
	GetName() *plugin.Kind
	GetDynamicResourceData(*resource.StaticResData) (*resource.DynamicResData, error)
}
