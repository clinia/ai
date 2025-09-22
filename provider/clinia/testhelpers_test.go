package clinia

import (
	"context"
	"errors"

	"github.com/clinia/models-client-go/cliniamodel/common"
)

// requesterStub implements common.Requester with minimal behaviour for unit tests.
type requesterStub struct{}

func (requesterStub) Infer(ctx context.Context, req common.InferRequest) (*common.InferResponse, error) {
	return nil, errors.New("not implemented")
}

func (requesterStub) Stream(ctx context.Context, modelName, modelVersion string, inputs []common.Input) (chan<- string, error) {
	return nil, errors.New("not implemented")
}

func (requesterStub) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func (requesterStub) Health(ctx context.Context) error { return nil }

func (requesterStub) Close() error { return nil }
