package aggservice

import (
	"context"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/go-kit/log"
)

const basePrice = 3.15

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store Storer
}

func (svc *BasicService) Aggregate(ctx context.Context, distance types.Distance) error {
	return svc.store.Insert(distance)
}

func (svc *BasicService) Calculate(ctx context.Context, obuId int) (*types.Invoice, error) {
	totalDistance, err := svc.store.Get(obuId)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		OBUID:         obuId,
		TotalDistance: totalDistance,
		TotalAmount:   totalDistance * basePrice,
	}

	return invoice, nil
}

func newBasicService(store Storer) Service {
	return &BasicService{store: store}
}

// NewAggregatorService will construct complete microservice
// with logging and instrumentation middleware
func New(logger log.Logger) Service {
	logger = log.With(logger, "service", "aggregator")

	var svc Service
	{
		svc = newBasicService(NewMemoryStore())
		svc = newLoggingMiddleware(logger)(svc)
		svc = newInstrumentationMiddleware()(svc)
	}

	return svc
}
