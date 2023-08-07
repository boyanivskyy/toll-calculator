package aggservice

import (
	"context"

	"github.com/boyanivskyy/toll-calculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func (mw loggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return nil
}

func (mw loggingMiddleware) Calculate(ctx context.Context, totalDistance int) (*types.Invoice, error) {
	return nil, nil
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next: next,
		}
	}
}

type instrumentationMiddleware struct {
	next Service
}

func (imw *instrumentationMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return nil
}

func (imw *instrumentationMiddleware) Calculate(ctx context.Context, totalDistance int) (*types.Invoice, error) {
	return nil, nil
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}
