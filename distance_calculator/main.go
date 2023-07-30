package main

import (
	"github.com/boyanivskyy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var kafkaTopic = "obudata"

const httpAggregatorEndpoint = "http://localhost:3000/aggregate"
const grpcAggEndpoint = "localhost:3001"

// Transport can be HTTP, gRPC, Kafka
// attach business logic to this transport

func main() {
	service := NewCalculatorService()
	service = NewLogMiddleware(service)

	// httpClient := client.NewHttpClient(httpAggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(grpcAggEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, service, grpcClient)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
