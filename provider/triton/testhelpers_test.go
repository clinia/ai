package triton

import (
	"context"
	"errors"

	"github.com/clinia/models-client-go/cliniamodel/common"
)

// requesterStub implements common.Requester with minimal behaviour for unit tests.
type requesterStub struct {
	inferCalls  int
	streamCalls int
	closeCalls  int
	closeErr    error
}

func (r *requesterStub) Infer(ctx context.Context, req common.InferRequest) (*common.InferResponse, error) {
	r.inferCalls++
	return nil, errors.New("not implemented")
}

func (r *requesterStub) Stream(ctx context.Context, modelName, modelVersion string, inputs []common.Input) (chan<- string, error) {
	r.streamCalls++
	return nil, errors.New("not implemented")
}

func (r *requesterStub) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func (r *requesterStub) Health(ctx context.Context) error { return nil }

func (r *requesterStub) Close() error {
	r.closeCalls++
	return r.closeErr
}
