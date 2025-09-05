package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	jinaprovider "go.jetify.com/ai/provider/jina"
	jina "go.jetify.com/ai/provider/jina/client"
)

func example() error {
	// Initialize the OpenAI provider
	provider := jinaprovider.NewProvider()

	// Create a model
	model := provider.NewMultimodalEmbeddingModel("jina-embeddings-v4")

	Text1 := "Artificial intelligence is the simulation of human intelligence in machines."

	Text2 := "Machine learning is a subset of AI that enables systems to learn from data."
	// Generate text
	task := "retrieval.query"
	response, err := ai.EmbedMany(
		context.Background(),
		model,
		[]jina.MultimodalEmbeddingInput{
			{Text: &Text1},
			{Text: &Text2},
		},
		ai.WithEmbeddingProviderMetadata[jina.MultimodalEmbeddingInput](
			"jina",
			jina.MultimodalEmbeddingNewParams{
				Task: &task,
			},
		),
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
