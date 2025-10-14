package main

import (
	"context"
	"log"
	"os"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina"
	jinaclient "go.jetify.com/ai/provider/jina/client"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Initialize the Jina provider (reads JINA_API_KEY/JINA_BASE_URL from env if set)
	provider := jina.NewProvider()

	// Create a Segmenter model (modelID used for metadata/logging only)
	model, err := provider.Segmenter("segmenter:1")
	if err != nil {
		return err
	}

	texts := []string{"Hello", "World"}

	// Enable true batching by passing provider metadata for Jina
	resp, err := ai.SegmentMany(
		ctx,
		model,
		texts,
		ai.WithTransportProviderMetadata("jina", jinaclient.SegmenterNewParams{UseContentArray: true}),
		ai.WithTransportBaseURL("https://model-7wl0980w.api.baseten.co/environments/production/predict"),
		ai.WithTransportAPIKey(os.Getenv("BASETEN_API_KEY")),
		ai.WithTransportUseRawBaseURL(),
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
