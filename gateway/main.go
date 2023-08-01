package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address of gateway server")
	flag.Parse()

	client := client.NewHttpClient("http://localhost:3000")
	invoiceHandler := NewInvoiceHandler(client)
	http.HandleFunc("/invoice", MakeApiFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port %s", *listenAddr)
	http.ListenAndServe(*listenAddr, nil)
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	invoice, err := h.client.GetInvoice(context.Background(), 1)
	if err != nil {
		// return writeJSON(w, http.StatusInternalServerError, map[string]string{
		// 	"error": err.Error(),
		// })
		return err
	}
	return writeJSON(w, http.StatusOK, invoice)
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}

func MakeApiFunc(fn ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}
}
