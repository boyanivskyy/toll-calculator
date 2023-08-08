package aggservice

import (
	"context"
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/go-kit/log"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	log  log.Logger
	next Service
}

func (mw loggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	err := mw.next.Aggregate(ctx, distance)

	defer func(start time.Time) {
		mw.log.Log("took", time.Since(start), "obuId", distance.OBUID, "value", distance.Value, "unix", distance.Unix, "err", err)
	}(time.Now())

	return err
}

func (mw loggingMiddleware) Calculate(ctx context.Context, obuId int) (*types.Invoice, error) {
	invoice, err := mw.next.Calculate(ctx, obuId)

	defer func(start time.Time) {
		mw.log.Log("took", time.Since(start), "obuId", invoice.OBUID, "totalDistance", invoice.TotalDistance, "totalAmount", invoice.TotalAmount, "err", err)
	}(time.Now())

	return invoice, err
}

func newLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next: next,
			log:  logger,
		}
	}
}

type instrumentationMiddleware struct {
	next Service
}

func (imw *instrumentationMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return imw.next.Aggregate(ctx, distance)
}

func (imw *instrumentationMiddleware) Calculate(ctx context.Context, obuId int) (*types.Invoice, error) {
	return imw.next.Calculate(ctx, obuId)
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}
