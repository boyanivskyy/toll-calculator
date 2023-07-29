package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/boyanivskyy/toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpListenaddr", ":3000", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcListenAddr", ":3001", "the listen address of the GRPC server")
	flag.Parse()

	store := NewMemoryStore()
	service := NewInvoiceAggregator(store)
	service = NewLoggingMiddleware(service)
	go makeGRPCTransport(*grpcListenAddr, service)
	makeHTTPTransport(*httpListenAddr, service)
}

func makeGRPCTransport(listenAddress string, service Aggregator) error {
	fmt.Println("GRPC transporter running on port", listenAddress)
	// Make a TCP listener
	listener, err := net.Listen("TCP", listenAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	// Make a new gRPC native server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register gRPC server implementation into gRPC server package
	types.RegisterAggregatorServer(server, NewGRPSServer(service))
	return server.Serve(listener)
}

func makeHTTPTransport(listenAddress string, service Aggregator) {
	fmt.Println("HTTP transporter running on port", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.HandleFunc("/invoice", handleGetInvoice(service))
	http.ListenAndServe(listenAddress, nil)
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obuId"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "No obuId",
			})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid of obuId",
			})
			return
		}

		invoice, err := service.CalculateInvoice(obuId)
		if err != nil {
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
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if err := service.AggregateDistance(distance); err != nil {
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
