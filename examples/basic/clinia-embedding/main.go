package main

import (
	"context"
	"log"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia"
)

func example() error {
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
	defer func() {
		_ = requester.Close()
	}()

	provider, err := clinia.NewProvider(ctx, clinia.WithClientOptions(common.ClientOptions{Requester: requester}))
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
	)
	if err != nil {
		return err
	}

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
