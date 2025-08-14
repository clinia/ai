package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp/v3"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	openaiprovider "go.jetify.com/ai/provider/openai"
)

func example() error {
	// Set up client options for the OpenAI client
	apiKey := os.Getenv("BASETEN_API_KEY")
	modelID := os.Getenv("BASETEN_MODEL_ID")
	modelEnv := os.Getenv("BASETEN_MODEL_ENV")

	clientOptions := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithMaxRetries(0), // Disable retries
	}

	// Create client with options
	client := openai.NewClient(clientOptions...)

	// Initialize the OpenAI provider
	provider := openaiprovider.NewProvider(openaiprovider.WithClient(client))

	// Create a model
	// NOTE: for baseten models we need to render the base URL in respect to the model ID
	model := provider.NewEmbeddingModel("text-embedding-3-small")

	baseURL := fmt.Sprintf("https://%s.api.baseten.co/environments/%s/sync/v1", modelID, modelEnv)

	// Generate text
	response, err := ai.EmbedMany(
		context.Background(),
		model,
		[]string{
			"Artificial intelligence is the simulation of human intelligence in machines.",
			"Machine learning is a subset of AI that enables systems to learn from data.",
		},
		ai.WithEmbeddingBaseURL[string](baseURL),
	)
	if err != nil {
		return err
	}

	// Print the response:
	printResponse(response)

	return nil
}

func printResponse(response api.EmbeddingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Print(response)
}

func main() {
	if err := example(); err != nil {
		log.Fatal(err)
	}
}
