package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	chonkie "go.jetify.com/ai/provider/chonkie"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Initialize the Chonkie provider (reads CHONKIE_API_KEY/CHONKIE_BASE_URL from env if set)
	provider := chonkie.NewProvider()

	// Create a Segmenting model (modelID used for metadata/logging only)
	model, err := provider.SegmentingModel("segmenting:1")
	if err != nil {
		return err
	}

	texts := []string{"Hello", "World"}
	godotenv.Load()

	// Enable true batching by passing provider metadata for Chonkie
	resp, err := ai.SegmentMany(
		ctx,
		model,
		texts,
		ai.WithTransportBaseURL(os.Getenv("BASETEN_SEGMENTING_URL")),
		ai.WithTransportAPIKey(os.Getenv("BASETEN_API_KEY")),
	)
	if err != nil {
		return err
	}

	printSegmentingResponse(texts, resp)
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
