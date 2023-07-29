package main

import (
	"fmt"

	"github.com/boyanivskyy/toll-calculator/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("Processing and inserting distance in the storage", distance)
	return i.store.Insert(distance)
}
func (i *InvoiceAggregator) CalculateInvoice(obuId int) (*types.Invoice, error) {
	totalDistance, err := i.store.Get(obuId)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		OBUID:         obuId,
		TotalDistance: totalDistance,
		TotalAmount:   totalDistance * basePrice,
	}

	return invoice, nil
}
