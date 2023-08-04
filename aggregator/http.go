package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPMetricsHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricsHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricsHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricsHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
			}).Info()
			h.reqLatency.Observe(latency)
		}(time.Now())

		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "method not supported",
			})
		}

		values, ok := r.URL.Query()["obuId"]
		if !ok {
			logrus.Error("/invoice: No obuId ")
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "No obuId",
			})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			logrus.Errorf("/invoice: Invalid obuId(%d)", obuId)
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid of obuId",
			})
			return
		}

		invoice, err := service.CalculateInvoice(obuId)
		if err != nil {
			logrus.Errorf("/invoice: Error calculating invoice(%s)", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "method not supported",
			})
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			logrus.Errorf("/aggregate: Error decoding request body(%s)", err)
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if err := service.AggregateDistance(distance); err != nil {
			logrus.Errorf("/aggregate: Error aggregating distance(%s)", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}
