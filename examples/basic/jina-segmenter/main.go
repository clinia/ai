package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina"
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

	texts := []string{
		"Jina AI: Your Search Foundation, Supercharged! 🚀\nIhrer Suchgrundlage, aufgeladen! 🚀\n您的搜索底座，从此不同！🚀\n検索ベース,もう二度と同じことはありません！🚀",
	}

	// Call the segmenter. Optionally set a custom base URL or headers:
	// resp, err := ai.SegmentMany(ctx, model, texts, ai.WithSegmentingBaseURL("https://api.jina.ai/v1/"))
	resp, err := ai.SegmentMany(ctx, model, texts)
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
