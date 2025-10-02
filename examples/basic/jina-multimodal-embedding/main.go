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
	// Initialize the OpenAI provider
	provider := jinaprovider.NewProvider()

	// Create a model
	model, _ := provider.MultimodalEmbeddingModel("jina-embeddings-v4")

	Text1 := "Artificial intelligence is the simulation of human intelligence in machines."

	Text2 := "Machine learning is a subset of AI that enables systems to learn from data."
	// Generate text
	task := "retrieval.query"
	response, err := ai.EmbedMany(
		context.Background(),
		model,
		[]api.MultimodalEmbeddingInput{
			{Text: &Text1},
			{Text: &Text2},
		},
		ai.WithEmbeddingProviderMetadata("jina", map[string]any{"task": task}),
	)
	if err != nil {
		return err
	}

	// Print the response:
	printResponse(response)

	return nil
}

func printResponse(response api.DenseEmbeddingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Print(response)
}

func main() {
	if err := example(); err != nil {
		log.Fatal(err)
	}
}
