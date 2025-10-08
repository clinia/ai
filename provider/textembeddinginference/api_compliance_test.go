package textembeddinginference

import (
	"encoding/json"
	"testing"

	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// TestEmbedRequestCompliance verifies our EmbedRequest structure matches TEI OpenAPI spec
func TestEmbedRequestCompliance(t *testing.T) {
	// Test that our EmbedRequest can marshal to the expected JSON structure
	req := tei.EmbedRequest{
		Inputs:              []string{"Hello world", "This is a test"},
		Dimensions:          intPtr(768),
		Normalize:           boolPtr(true),
		Truncate:            boolPtr(false),
		TruncationDirection: stringPtr("Right"),
		PromptName:          stringPtr("query"),
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal EmbedRequest: %v", err)
	}

	// Verify JSON structure matches expected format
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required field
	if _, ok := decoded["inputs"]; !ok {
		t.Error("Missing required 'inputs' field")
	}

	// Check optional fields are present when set
	expectedFields := []string{"dimensions", "normalize", "truncate", "truncation_direction", "prompt_name"}
	for _, field := range expectedFields {
		if _, ok := decoded[field]; !ok {
			t.Errorf("Missing expected field: %s", field)
		}
	}

	t.Logf("EmbedRequest JSON: %s", string(jsonData))
}

// TestEmbedSparseRequestCompliance verifies our EmbedSparseRequest structure matches TEI OpenAPI spec
func TestEmbedSparseRequestCompliance(t *testing.T) {
	req := tei.EmbedSparseRequest{
		Inputs:              []string{"Hello world", "This is a test"},
		Truncate:            boolPtr(false),
		TruncationDirection: stringPtr("Right"),
		PromptName:          stringPtr("query"),
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal EmbedSparseRequest: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required field
	if _, ok := decoded["inputs"]; !ok {
		t.Error("Missing required 'inputs' field")
	}

	t.Logf("EmbedSparseRequest JSON: %s", string(jsonData))
}

// TestSparseValueCompliance verifies our SparseValue structure matches TEI OpenAPI spec
func TestSparseValueCompliance(t *testing.T) {
	sv := tei.SparseValue{
		Index: 42,
		Value: 0.75,
	}

	jsonData, err := json.Marshal(sv)
	if err != nil {
		t.Fatalf("Failed to marshal SparseValue: %v", err)
	}

	// Verify it matches expected format: {"index": 42, "value": 0.75}
	expected := `{"index":42,"value":0.75}`
	if string(jsonData) != expected {
		t.Errorf("SparseValue JSON mismatch. Expected: %s, Got: %s", expected, string(jsonData))
	}

	t.Logf("SparseValue JSON: %s", string(jsonData))
}

// TestResponseStructures verifies our response types can handle TEI API responses
func TestResponseStructures(t *testing.T) {
	// Test dense embedding response (should be [][]float64)
	denseResponse := tei.CreateEmbeddingResponse{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
	}

	denseJSON, err := json.Marshal(denseResponse)
	if err != nil {
		t.Fatalf("Failed to marshal dense response: %v", err)
	}
	t.Logf("Dense embedding response JSON: %s", string(denseJSON))

	// Test sparse embedding response
	sparseResponse := tei.CreateSparseEmbeddingResponse{
		{
			{Index: 10, Value: 0.5},
			{Index: 20, Value: 0.8},
		},
		{
			{Index: 15, Value: 0.3},
			{Index: 25, Value: 0.9},
		},
	}

	sparseJSON, err := json.Marshal(sparseResponse)
	if err != nil {
		t.Fatalf("Failed to marshal sparse response: %v", err)
	}
	t.Logf("Sparse embedding response JSON: %s", string(sparseJSON))
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
