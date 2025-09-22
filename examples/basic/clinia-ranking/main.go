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

	ranker, err := provider.NewRankingModel("hybrid-ranker:1")
	if err != nil {
		return err
	}

	query := "Is the company based in Montreal?"
	texts := []string{
		"Clinia is based in Montreal",
		"The other facility is located in Toronto",
		"That house was built in 1990",
		"Neuroscience is the study of the nervous system",
		"The conference will be held in Paris next year",
		"The new product will be launched in Q3 2022",
		"One of the largest cities in the world is Tokyo",
		"Before the meeting, please read the document.",
		"The patient was admitted to the hospital yesterday",
		"The event will take place at the convention center",
	}

	resp, err := ai.RankMany(ctx, ranker, query, texts)
	if err != nil {
		return err
	}

	printRankingResponse(query, resp)
	return nil
}

func printRankingResponse(query string, resp api.RankingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Println(map[string]any{
		"query":  query,
		"scores": resp.Scores,
	})
}
