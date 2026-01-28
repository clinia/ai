package instrumentation

import "context"

// Attributes are optional key/value pairs attached to spans.
// They are intentionally untyped so callers can adapt them to their tracing backend.
type Attributes map[string]any

// Operation describes the kind of provider call being traced.
type Operation string

const (
	OperationGenerate Operation = "generate"
	OperationStream   Operation = "stream"
	OperationEmbed    Operation = "embed"
	OperationRank     Operation = "rank"
	OperationSegment  Operation = "segment"
)

// ProviderSpanInfo carries provider-specific identifiers for a span.
type ProviderSpanInfo struct {
	Provider  string
	Model     string
	Operation Operation
}

// Span is the handle returned from Instrumenter.Start.
type Span interface {
	End(err error)
}

// Instrumenter is implemented by tracing backends to create spans.
type Instrumenter interface {
	Start(ctx context.Context, spanName string, attributes Attributes, info ProviderSpanInfo) (context.Context, Span)
}

type noopInstrumenter struct{}
type noopSpan struct{}

// NopInstrumenter returns an Instrumenter that performs no tracing.
func NopInstrumenter() Instrumenter { return noopInstrumenter{} }

// Start implements Instrumenter for the noop instrumenter.
func (noopInstrumenter) Start(ctx context.Context, _ string, _ Attributes, _ ProviderSpanInfo) (context.Context, Span) {
	return ctx, noopSpan{}
}

// End implements Span for the noop span.
func (noopSpan) End(error) {}

// EndSpan ends the provided span with the given error pointer.
// It safely handles nil spans and nil error pointers.
func EndSpan(span Span, err *error) {
	if span == nil {
		return
	}

	var endErr error
	if err != nil {
		endErr = *err
	}

	span.End(endErr)
}
