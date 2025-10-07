package main

import (
	"context"
	"fmt"
	"log"

	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

func main() {
	fmt.Println("TEI Reranking Example")
	fmt.Println("====================")

	// Initialize TEI client - set TEI_BASE_URL environment variable
	// For reranking models like BAAI/bge-reranker-large
	// Example: export TEI_BASE_URL=http://localhost:8080
	client := tei.NewClient()

	ctx := context.Background()

	// Example query and documents to rerank
	query := "What is machine learning?"
	documents := []string{
		"Machine learning is a subset of artificial intelligence that uses algorithms to learn patterns from data.",
		"The weather forecast shows rain tomorrow with a 70% chance of precipitation.",
		"Deep learning uses neural networks with multiple layers to process complex data patterns.",
		"Cooking pasta requires boiling water and adding salt for better flavor.",
		"Supervised learning trains models using labeled examples to make predictions on new data.",
		"Basketball is a sport played with two teams of five players each on a rectangular court.",
		"Natural language processing helps computers understand and generate human language.",
		"The stock market experienced volatility due to economic uncertainty and inflation concerns.",
	}

	fmt.Printf("Query: %s\n\n", query)
	fmt.Printf("Documents to rerank (%d total):\n", len(documents))
	for i, doc := range documents {
		fmt.Printf("  [%d] %s\n", i, doc)
	}
	fmt.Println()

	// Example 1: Basic reranking
	fmt.Println("=== Example 1: Basic Reranking ===")
	basicRequest := tei.RerankRequest{
		Query: query,
		Texts: documents,
	}

	response, err := client.Reranking.Rerank(ctx, basicRequest)
	if err != nil {
		log.Fatalf("Basic reranking failed: %v", err)
	}

	fmt.Println("Results (sorted by relevance):")
	fmt.Println("Rank | Score | Index | Relevance")
	fmt.Println("-----|-------|-------|----------")
	for rank, result := range *response {
		relevance := "High"
		if result.Score < 0.5 {
			relevance = "Medium"
		}
		if result.Score < 0.2 {
			relevance = "Low"
		}
		fmt.Printf("%-4d | %.4f | %-5d | %s\n",
			rank+1, result.Score, result.Index, relevance)
	}
	fmt.Println()

	// Example 2: Reranking with text return and options
	fmt.Println("=== Example 2: Reranking with Options ===")
	returnText := true
	truncate := false
	rawScores := false

	optionsRequest := tei.RerankRequest{
		Query:      query,
		Texts:      documents[:5], // Use first 5 documents for cleaner output
		ReturnText: &returnText,
		RawScores:  &rawScores,
		Truncate:   &truncate,
	}

	response2, err := client.Reranking.Rerank(ctx, optionsRequest)
	if err != nil {
		log.Fatalf("Options reranking failed: %v", err)
	}

	fmt.Println("Detailed Results with Text:")
	fmt.Println("Rank | Score | Index | Document Text")
	fmt.Println("-----|-------|-------|-------------")
	for rank, result := range *response2 {
		text := "N/A"
		if result.Text != nil {
			// Truncate text for display
			if len(*result.Text) > 80 {
				text = (*result.Text)[:77] + "..."
			} else {
				text = *result.Text
			}
		}
		fmt.Printf("%-4d | %.4f | %-5d | %s\n",
			rank+1, result.Score, result.Index, text)
	}
	fmt.Println()

	// Example 3: Reranking with truncation settings
	fmt.Println("=== Example 3: Reranking with Truncation ===")
	truncateEnabled := true
	truncationDir := "Right"

	truncateRequest := tei.RerankRequest{
		Query:               query,
		Texts:               documents[:3], // Use first 3 for demo
		Truncate:            &truncateEnabled,
		TruncationDirection: &truncationDir,
		ReturnText:          &returnText,
	}

	response3, err := client.Reranking.Rerank(ctx, truncateRequest)
	if err != nil {
		log.Fatalf("Truncation reranking failed: %v", err)
	}

	fmt.Printf("Results with truncation (direction: %s):\n", truncationDir)
	for rank, result := range *response3 {
		fmt.Printf("Rank %d: Score %.4f, Original Index %d\n",
			rank+1, result.Score, result.Index)
		if result.Text != nil {
			fmt.Printf("  Text: %s\n", *result.Text)
		}
	}
	fmt.Println()

	// Summary statistics
	fmt.Println("=== Summary ===")
	fmt.Printf("Total documents processed: %d\n", len(documents))
	fmt.Printf("Highest relevance score: %.4f\n", (*response)[0].Score)
	fmt.Printf("Most relevant document index: %d\n", (*response)[0].Index)
	fmt.Printf("Most relevant document: %s\n", documents[(*response)[0].Index])

	// Count high relevance documents (score > 0.5)
	highRelevance := 0
	for _, result := range *response {
		if result.Score > 0.5 {
			highRelevance++
		}
	}
	fmt.Printf("Documents with high relevance (>0.5): %d\n", highRelevance)
}
