package commands

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var (
	apiClient *http.Client
	baseUrl   string
)

func init() {
	//Timeout measures the maximum time a request can take before it is considered failed.
	// Transport is the HTTP transport used by the client to send requests.
	// TLSClientConfig is the configuration for TLS (Transport Layer Security) connections.
	// InsecureSkipVerify is a boolean that indicates whether to skip the verification of the server's TLS certificate.
	// This is useful for development or testing purposes, but should not be used in production.
	// InsecureSkipVerify is set to true to ignore TLS certificate verification.
	// This is not recommended for production use, as it can expose you to security risks.
	apiClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Ignora a verificação de certificado TLS
			},
		},
	}
	// viper é uma biblioteca de configuração que permite ler e escrever configurações de forma fácil.
	// SetDefault é usado para definir um valor padrão para uma chave de configuração.
	viper.SetDefault("api.url", "http://localhost:8080") // Define a URL base do servidor API
	//GetString é usado para obter o valor de uma chave de configuração como uma string.
	baseUrl = viper.GetString("api.url") // URL base do servidor API
}
