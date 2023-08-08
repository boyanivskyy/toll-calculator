package aggendpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(start time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(start).Seconds())
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(start time.Time) {
				logger.Log("error", err, "took", time.Since(start))
			}(time.Now())

			return next(ctx, request)
		}
	}
}
