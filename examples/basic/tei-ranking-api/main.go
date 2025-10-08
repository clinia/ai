package main

import (
	"context"
	"fmt"
	"log"

	"go.jetify.com/ai"
	tei "go.jetify.com/ai/provider/textembeddinginference"
	teiclient "go.jetify.com/ai/provider/textembeddinginference/client"
)

func main() {
	fmt.Println("TEI Ranking API Compliance Example")
	fmt.Println("==================================")

	// Initialize TEI provider
	// Set TEI_BASE_URL environment variable to your TEI reranking service
	// Example: export TEI_BASE_URL=http://localhost:8080
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

	// Use TEI-specific options through the API
	teiOpts := ai.TEIRankingOptions{
		ReturnText:          boolPtr(true),
		RawScores:           boolPtr(false),
		Truncate:            boolPtr(false),
		TruncationDirection: stringPtr("Right"),
	}

	response2, err := ai.RankMany(ctx, rankingModel, query, documents, ai.WithTEIRankingOptions(teiOpts))
	if err != nil {
		log.Fatalf("TEI-specific ranking failed: %v", err)
	}

	fmt.Println("API Response with TEI options:")
	fmt.Println("Index | Score")
	fmt.Println("------|------")
	for i, score := range response2.Scores {
		fmt.Printf("%-5d | %.4f\n", i, score)
	}
	fmt.Println()

	// Example 3: Direct TEI client usage for detailed results
	fmt.Println("=== Example 3: Direct TEI Client for Detailed Results ===")

	// For detailed results with indices and text, use the TEI client directly
	client := teiclient.NewClient()

	detailedRequest := teiclient.RankRequest{
		Query:      query,
		Texts:      documents[:3], // Use fewer documents for cleaner output
		ReturnText: boolPtr(true),
		RawScores:  boolPtr(false),
		Truncate:   boolPtr(false),
	}

	detailedResp, err := client.Ranking.Rank(ctx, detailedRequest)
	if err != nil {
		log.Fatalf("Direct TEI ranking failed: %v", err)
	}

	fmt.Println("Detailed TEI Response (with indices and text):")
	fmt.Println("Rank | Score | Index | Text")
	fmt.Println("-----|-------|-------|-----")
	for rank, result := range *detailedResp {
		text := "N/A"
		if result.Text != nil {
			if len(*result.Text) > 50 {
				text = (*result.Text)[:47] + "..."
			} else {
				text = *result.Text
			}
		}
		fmt.Printf("%-4d | %.4f | %-5d | %s\n", rank+1, result.Score, result.Index, text)
	}

	// Comparison
	fmt.Printf("\nComparison:\n")
	fmt.Printf("Standard API: Returns scores in original document order\n")
	fmt.Printf("TEI Direct: Returns results sorted by relevance with indices\n")

	fmt.Printf("\nModel Info:\n")
	fmt.Printf("Provider: %s\n", rankingModel.ProviderName())
	fmt.Printf("Model ID: %s\n", rankingModel.ModelID())
	fmt.Printf("Specification Version: %s\n", rankingModel.SpecificationVersion())
	fmt.Printf("Supports Parallel Calls: %v\n", rankingModel.SupportsParallelCalls())
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
