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

	segmenter, err := provider.Segmenter("text-chunker:1")
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

	resp, err := ai.SegmentMany(ctx, segmenter, documents, ai.WithSegmentingBaseURL("http://127.0.0.1:4770"))
	if err != nil {
		return err
	}

	printSegmentingResponse(documents, resp)
	return nil
}

func printSegmentingResponse(texts []string, resp api.SegmentingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Println(map[string]any{
		"texts":    texts,
		"segments": resp.Segments,
	})
}
