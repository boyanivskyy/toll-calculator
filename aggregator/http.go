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

type ApiError struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

// Error implements the Error interface
func (e ApiError) Error() string {
	return e.Err.Error()
}

type HttpFunc func(w http.ResponseWriter, r *http.Request) error

type HTTPMetricsHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHttpHandlerFunc(fn HttpFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(ApiError); ok {
				writeJSON(w, apiErr.Code, apiErr)
			}
		}
	}
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricsHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "error_counter"),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricsHandler{
		reqCounter: reqCounter,
		errCounter: errCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricsHandler) instrument(next HttpFunc) HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err":     err,
			}).Info()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		return err
	}
}

func handleGetInvoice(service Aggregator) HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("method not supported, expect %s, got %s", "GET", r.Method),
			}
		}

		values, ok := r.URL.Query()["obuId"]
		if !ok {
			logrus.Error("/invoice: No obuId ")
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("no obuID"),
			}
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			logrus.Errorf("/invoice: Invalid obuId(%d)", obuId)
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid obuID"),
			}
		}

		invoice, err := service.CalculateInvoice(obuId)
		if err != nil {
			logrus.Errorf("/invoice: Error calculating invoice(%s)", err)
			return ApiError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(service Aggregator) HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("method not supported, handling POST, get %s", r.Method),
			}
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			logrus.Errorf("/aggregate: Error decoding request body(%s)", err)
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("Error decoding request body %s", err),
			}
		}
		if err := service.AggregateDistance(distance); err != nil {
			logrus.Errorf("/aggregate: Error aggregating distance: (%s)", err)
			return ApiError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("Error aggregating distance: %s", err),
			}
		}

		return writeJSON(w, http.StatusOK, map[string]string{
			"message": "ok",
		})
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
