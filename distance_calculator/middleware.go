package main

import (
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(service CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: service,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (float64, error) {
	var (
		dist float64
		err  error
	)
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": dist,
		}).Info("calculate distance")
	}(time.Now())

	dist, err = m.next.CalculateDistance(data)
	return dist, err
}
