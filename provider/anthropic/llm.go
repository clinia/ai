package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	"go.jetify.com/ai/provider/anthropic/codec"
)

// ModelOption is a function type that modifies a LanguageModel.
type ModelOption func(*LanguageModel)

// WithClient returns a ModelOption that sets the client.
func WithClient(client anthropic.Client) ModelOption {
	// TODO: Instead of only supporting an anthropic.Client, we can "flatten"
	// the options supported by the Anthropic SDK.
	return func(m *LanguageModel) {
		m.client = client
	}
}

// WithInstrumenter configures tracing for model calls.
func WithInstrumenter(instr instrumentation.Instrumenter) ModelOption {
	return func(m *LanguageModel) {
		if instr == nil {
			instr = instrumentation.NopInstrumenter()
		}
		m.instrumenter = instr
	}
}

// LanguageModel represents an Anthropic language model.
type LanguageModel struct {
	modelID      string
	client       anthropic.Client
	instrumenter instrumentation.Instrumenter
}

var _ api.LanguageModel = &LanguageModel{}

// NewLanguageModel creates a new Anthropic language model.
func NewLanguageModel(modelID string, opts ...ModelOption) *LanguageModel {
	// Create model with default settings
	model := &LanguageModel{
		modelID:      modelID,
		client:       anthropic.NewClient(), // Default client
		instrumenter: instrumentation.NopInstrumenter(),
	}

	// Apply options
	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (m *LanguageModel) ProviderName() string {
	return ProviderName
}

func (m *LanguageModel) ModelID() string {
	return m.modelID
}

func (m *LanguageModel) SupportedUrls() []api.SupportedURL {
	// TODO: Make configurable via the constructor.
	return []api.SupportedURL{
		{
			MediaType: "image/*",
			URLPatterns: []string{
				"^https?://.*",
			},
		},
	}
}

func (m *LanguageModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (resp *api.Response, err error) {
	ctx, span := m.instrumenter.Start(
		ctx,
		"Generate",
		instrumentation.Attributes{
			"provider":   m.ProviderName(),
			"model":      m.modelID,
			"model_type": "language_model",
			"operation":  string(instrumentation.OperationGenerate),
		},
		instrumentation.ProviderSpanInfo{
			Provider:  m.ProviderName(),
			Model:     m.modelID,
			Operation: instrumentation.OperationGenerate,
		},
	)
	defer instrumentation.EndSpan(span, &err)

	params, warnings, err := codec.EncodeParams(m.modelID, prompt, opts)
	if err != nil {
		return nil, err
	}

	message, err := m.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	response, err := codec.DecodeResponse(message)
	if err != nil {
		return nil, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (resp *api.StreamResponse, err error) {
	ctx, span := m.instrumenter.Start(
		ctx,
		"Stream",
		instrumentation.Attributes{
			"provider":   m.ProviderName(),
			"model":      m.modelID,
			"model_type": "language_model",
			"operation":  string(instrumentation.OperationStream),
		},
		instrumentation.ProviderSpanInfo{
			Provider:  m.ProviderName(),
			Model:     m.modelID,
			Operation: instrumentation.OperationStream,
		},
	)
	defer instrumentation.EndSpan(span, &err)

	return nil, api.NewUnsupportedFunctionalityError("streaming generation", "")
}
