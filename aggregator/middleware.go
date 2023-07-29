package main

import (
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

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
