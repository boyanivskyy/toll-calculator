package main

import (
	"net"
	"net/http"
	"os"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	grpcListenAddr := os.Getenv("AGG_GRPC_ENDPOINT")
	httpListenAddr := os.Getenv("AGG_HTTP_ENDPOINT")
	store := makeStore()
	service := NewInvoiceAggregator(store)
	service = NewMetricsMiddleware(service)
	service = NewLoggingMiddleware(service)

	go func() {
		if err := makeGRPCTransport(grpcListenAddr, service); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Fatal(makeHTTPTransport(httpListenAddr, service))
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
	var (
		aggregateMetricsHandler = newHTTPMetricsHandler("aggregate")
		invoiceMetricsHandler   = newHTTPMetricsHandler("invoice")
		invoiceHandler          = makeHttpHandlerFunc(invoiceMetricsHandler.instrument(handleGetInvoice(service)))
		aggregateHandler        = makeHttpHandlerFunc(aggregateMetricsHandler.instrument(handleAggregate(service)))
	)

	http.HandleFunc("/aggregate", aggregateHandler)
	http.HandleFunc("/invoice", invoiceHandler)
	http.Handle("/metrics", promhttp.Handler())

	logrus.Info("HTTP transporter running on port", listenAddress)
	return http.ListenAndServe(listenAddress, nil)
}

func makeStore() Storer {
	t := os.Getenv("AGG_STORE_TYPE")
	switch t {
	case "memory":
		return NewMemoryStore()
	default:
		logrus.Fatalf("invalid store type %s", t)
		return nil
	}
}
