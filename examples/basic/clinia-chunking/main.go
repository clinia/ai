package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	clinia "go.jetify.com/ai/provider/clinia"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	provider, err := clinia.NewProvider(ctx)
	if err != nil {
		return err
	}

	chunker, err := provider.ChunkingModel("text-chunker:1")
	if err != nil {
		return err
	}

	documents := []string{
		"Artificial intelligence is revolutionizing healthcare by enabling faster diagnoses and personalized treatment plans.",
		"Recent advancements in natural language processing allow AI models to understand context and generate human-like responses.",
		"Medical research journals play a crucial role in disseminating new findings, ensuring that professionals stay up to date.",
		"Deep learning techniques have significantly improved image recognition in radiology, assisting doctors in detecting anomalies early.",
		"Ethical considerations in AI development focus on transparency, fairness, and bias mitigation to ensure equitable outcomes.",
	}

	resp, err := ai.ChunkMany(ctx, chunker, documents, ai.WithChunkingBaseURL("http://127.0.0.1:4770"))
	if err != nil {
		return err
	}

	printChunkingResponse(documents, resp)
	return nil
}

func printChunkingResponse(texts []string, resp api.ChunkingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Println(map[string]any{
		"texts":  texts,
		"chunks": resp.Chunks,
	})
}
