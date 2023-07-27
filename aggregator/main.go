package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/types"
)

func main() {
	listenAddress := flag.String("listenaddr", ":3000", "the listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	service := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddress, service)
}

func makeHTTPTransport(listenAddress string, service Aggregator) {
	fmt.Println("HTTP transporter running on port", listenAddress)

	http.HandleFunc("/aggregate", handleAggregate(service))
	http.ListenAndServe(listenAddress, nil)
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// if err := service.AggregateDistance(distance); err != nil {
		// 	w.Write(json.Unmarshal(distance))
		// }
	}
}
