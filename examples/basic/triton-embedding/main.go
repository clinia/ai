package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/triton"
)

func example() error {
	ctx := context.Background()

	provider, err := triton.NewProvider()
	if err != nil {
		return err
	}

	model, err := provider.TextEmbeddingModel("dense-embedder:1")
	if err != nil {
		return err
	}

	response, err := ai.EmbedMany(
		ctx,
		model,
		[]string{
			"Hello, how are you?",
		},
		ai.WithTransportBaseURL("http://127.0.0.1:4770"),
	)
	if err != nil {
		return err
	}

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
