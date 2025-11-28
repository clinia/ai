package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/tei"
)

func main() {
	fmt.Println("TEI Sparse Embedding Example")
	fmt.Println("============================")

	// Example: export TEI_BASE_URL=http://localhost:8080

	// Load environment variables from .env if present
	if err := godotenv.Load(); err != nil {
		log.Printf("(.env load skipped) %v", err)
	}
	if os.Getenv("TEI_SPARSE_EMBEDDING_URL") == "" {
		log.Println("TEI_SPARSE_EMBEDDING_URL not set; set it in .env or your shell environment.")
	}

	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	provider := tei.NewProvider()

	// The model ID is used for logging/metadata; your TEI server decides which model runs.
	model, err := provider.SparseEmbeddingModel("BAAI/bge-base-en-v1.5")
	if err != nil {
		return err
	}

	texts := []string{
		"Artificial intelligence is transforming many industries.",
		"Baseball players train during the off-season to stay sharp.",
		"A neural network can learn complex relationships in data.",
	}

	fmt.Println("=== Example 1: Basic Sparse Embeddings ===")
	resp, err := ai.EmbedMany(ctx, model, texts, ai.WithTransportBaseURL(os.Getenv("TEI_SPARSE_EMBEDDING_URL")), ai.WithTransportAPIKey(os.Getenv("TEI_API_KEY")))
	if err != nil {
		return err
	}
	printSparseEmbeddings(texts, resp, 8)

	return nil
}

func printSparseEmbeddings(texts []string, resp api.SparseEmbeddingResponse, topK int) {
	if len(resp.Embeddings) != len(texts) {
		log.Printf("warning: got %d embeddings for %d inputs", len(resp.Embeddings), len(texts))
	}

	for i, emb := range resp.Embeddings {
		fmt.Printf("\nInput %d: %s\n", i+1, texts[i])
		fmt.Printf("Top %d weighted terms:\n", topK)
		for _, term := range topWeightedTerms(emb, topK) {
			fmt.Printf("  %s: %.4f\n", term.term, term.weight)
		}
	}
}

type weightedTerm struct {
	term   string
	weight float64
}

func topWeightedTerms(embedding api.SparseEmbedding, k int) []weightedTerm {
	terms := make([]weightedTerm, 0, len(embedding))
	for token, weight := range embedding {
		terms = append(terms, weightedTerm{
			term:   token,
			weight: weight,
		})
	}

	sort.Slice(terms, func(i, j int) bool {
		return terms[i].weight > terms[j].weight
	})

	if k > len(terms) {
		k = len(terms)
	}
	return terms[:k]
}
