package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.jetify.com/ai"
	"go.jetify.com/ai/provider/tei"
)

func main() {
	fmt.Println("TEI Ranking API Compliance Example")
	fmt.Println("==================================")

	// Example: export TEI_BASE_URL=http://localhost:8080

	// Load environment variables from .env if present
	if err := godotenv.Load(); err != nil {
		log.Printf("(.env load skipped) %v", err)
	}
	if os.Getenv("TEI_BASE_URL") == "" {
		log.Println("TEI_BASE_URL not set; set it in .env or your shell environment.")
	}

	provider := tei.NewProvider()
	ctx := context.Background()

	// Sample data

	// Example query and documents
	query := "What is machine learning?"
	documents := []string{
		"Machine learning is a subset of artificial intelligence that uses algorithms to learn patterns from data.",
		"The weather forecast shows rain tomorrow with a 70% chance of precipitation.",
		"Deep learning uses neural networks with multiple layers to process complex data patterns.",
		"Supervised learning trains models using labeled examples to make predictions on new data.",
		"Natural language processing helps computers understand and generate human language.",
	}

	fmt.Printf("Query: %s\n\n", query)
	fmt.Printf("Documents (%d total):\n", len(documents))
	for i, doc := range documents {
		fmt.Printf("  [%d] %s\n", i, doc)
	}
	fmt.Println()

	// Example 1: Using the standard API-compliant RankingModel interface
	fmt.Println("=== Example 1: API-Compliant Ranking (DoRank) ===")

	rankingModel, err := provider.RankingModel("BAAI/bge-reranker-large")
	if err != nil {
		log.Fatalf("Failed to create ranking model: %v", err)
	}

	// Use the standard ranking API (returns just scores)
	response, err := ai.RankMany(ctx, rankingModel, query, documents)
	if err != nil {
		log.Fatalf("Ranking failed: %v", err)
	}

	fmt.Println("API Response (scores only):")
	fmt.Println("Index | Score")
	fmt.Println("------|------")
	for i, score := range response.Scores {
		fmt.Printf("%-5d | %.4f\n", i, score)
	}
	fmt.Println()

	// Example 2: Using TEI-specific options with the API
	fmt.Println("=== Example 2: API with TEI-Specific Options ===")

	// Use TEI-specific options through ProviderMetadata
	teiProviderMetadata := map[string]any{
		"return_text":          true,
		"raw_scores":           false,
		"truncate":             false,
		"truncation_direction": "Right",
	}

	response2, err := ai.RankMany(ctx, rankingModel, query, documents, ai.WithTransportProviderMetadata("tei", teiProviderMetadata))
	if err != nil {
		log.Fatalf("TEI-specific ranking failed: %v", err)
	}

	fmt.Println("API Response with TEI options:")
	fmt.Println("Index | Score")
	fmt.Println("------|------")
	for i, score := range response2.Scores {
		fmt.Printf("%-5d | %.4f\n", i, score)
	}
}
