package main

import "github.com/sirupsen/logrus"

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var kafkaTopic = "obudata"

// Transport can be HTTP, gRPC, Kafka
// attach business logic to this transport

func main() {
	service := NewCalculatorService()
	service = NewLogMiddleware(service)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, service)
	if err != nil {
		logrus.Fatal(err)
	}

	kafkaConsumer.Start()
}
