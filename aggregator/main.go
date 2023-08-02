package main

import (
	"context"
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/boyanivskyy/toll-calculator/aggregator/client"
	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpListenaddr", ":3000", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcListenAddr", ":3001", "the listen address of the GRPC server")
	flag.Parse()

	store := NewMemoryStore()
	service := NewInvoiceAggregator(store)
	service = NewMetricsMiddleware(service)
	service = NewLoggingMiddleware(service)
	go func() {
		if err := makeGRPCTransport(*grpcListenAddr, service); err != nil {
			logrus.Fatal(err)
		}
	}()
	time.Sleep(time.Second * 2)
	c, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		logrus.Fatal("NewGRPCClient", err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 1,
		Value: 33.33,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		logrus.Fatal(err)
	}

	logrus.Fatal(makeHTTPTransport(*httpListenAddr, service))
}

func makeGRPCTransport(listenAddress string, service Aggregator) error {
	logrus.Info("GRPC transporter running on port", listenAddress)
	// Make a TCP listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	// Make a new gRPC native server
	server := grpc.NewServer(grpc.EmptyServerOption{})
	// Register gRPC server implementation into gRPC server package
	types.RegisterAggregatorServer(server, NewGRPSServer(service))
	return server.Serve(listener)
}

func makeHTTPTransport(listenAddress string, service Aggregator) error {
	logrus.Info("HTTP transporter running on port", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.HandleFunc("/invoice", handleGetInvoice(service))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAddress, nil)
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
