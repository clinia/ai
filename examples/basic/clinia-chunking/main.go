package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/triton"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	provider, err := triton.NewProvider(ctx)
	if err != nil {
		return err
	}

	segmenting, err := provider.SegmentingModel("text-chunker:1")
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

	resp, err := ai.SegmentMany(ctx, segmenting, documents, ai.WithTransportBaseURL("http://127.0.0.1:4770"))
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
