package jina

import "net/http"

type ProviderConfig struct {
	providerName string
	client       http.Client
	apiKey       string
}
