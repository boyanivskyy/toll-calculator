package main

import (
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterCalculateInvoice prometheus.Counter
	reqLatencyCalculateInvoice prometheus.Histogram
	errCounterCalculateInvoice prometheus.Counter

	reqCounterAggregate prometheus.Counter
	reqLatencyAggregate prometheus.Histogram
	errCounterAggregate prometheus.Counter

	next Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate_invoice",
	})

	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate_invoice",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate_invoice",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricsMiddleware{
		reqCounterAggregate:        reqCounterAgg,
		reqCounterCalculateInvoice: reqCounterCalc,

		reqLatencyAggregate:        reqLatencyAgg,
		reqLatencyCalculateInvoice: reqLatencyCalc,

		errCounterAggregate:        errCounterAgg,
		errCounterCalculateInvoice: errCounterCalc,

		next: next,
	}
}

func (m *MetricsMiddleware) AggregateDistance(data types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAggregate.Observe(time.Since(start).Seconds())
		m.reqCounterAggregate.Inc()
		if err != nil {
			m.errCounterAggregate.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(data)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuId int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalculateInvoice.Observe(time.Since(start).Seconds())
		m.reqCounterCalculateInvoice.Inc()
		if err != nil {
			m.errCounterCalculateInvoice.Inc()
		}
	}(time.Now())

	invoice, err = m.next.CalculateInvoice(obuId)
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLoggingMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(data types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())

	err = l.next.AggregateDistance(data)
	return err
}

func (l *LogMiddleware) CalculateInvoice(obuId int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		fields := logrus.Fields{
			"took":  time.Since(start),
			"err":   err,
			"obuId": obuId,
		}
		if invoice != nil {
			fields["totalDistance"] = invoice.TotalDistance
			fields["totalAmount"] = invoice.TotalAmount
		}

		logrus.WithFields(fields).Info("CalculateInvoice")
	}(time.Now())

	invoice, err = l.next.CalculateInvoice(obuId)
	return invoice, err
}
