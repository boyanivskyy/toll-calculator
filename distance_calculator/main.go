package main

import (
	"github.com/boyanivskyy/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var kafkaTopic = "obudata"

const aggregatorEndpoint = "http://localhost:3000/aggregate"

// Transport can be HTTP, gRPC, Kafka
// attach business logic to this transport

func main() {
	service := NewCalculatorService()
	service = NewLogMiddleware(service)
	client := client.NewHttpClient(aggregatorEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, service, client)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
