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

	ranker, err := provider.RankingModel("hybrid-ranker:1")
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

	resp, err := ai.RankMany(ctx, ranker, query, texts, ai.WithTransportBaseURL("http://127.0.0.1:4770"))
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
