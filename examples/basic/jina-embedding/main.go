package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	jinaprovider "go.jetify.com/ai/provider/jina"
)

func example() error {
	// Initialize the Jina provider
	provider := jinaprovider.NewProvider()

	// Create a model
	model, _ := provider.TextEmbeddingModel("jina-embeddings-v3")

	// Generate text
	response, err := ai.EmbedMany(
		context.Background(),
		model,
		[]string{
			"Artificial intelligence is the simulation of human intelligence in machines.",
			"Machine learning is a subset of AI that enables systems to learn from data.",
		},
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
