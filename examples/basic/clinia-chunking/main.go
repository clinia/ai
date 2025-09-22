package main

import (
	"context"
	"log"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
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

	host := common.Host{
		Url:    "127.0.0.1",
		Port:   4770,
		Scheme: common.HTTP,
	}

	requester, err := requestergrpc.NewRequester(ctx, common.RequesterConfig{Host: host})
	if err != nil {
		return err
	}
	defer requester.Close()

	provider, err := clinia.NewProvider(ctx, clinia.WithClientOptions(common.ClientOptions{Requester: requester}))
	if err != nil {
		return err
	}

	chunker, err := provider.NewChunkingModel("text-chunker", "1")
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

	resp, err := ai.ChunkMany(ctx, chunker, documents)
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
