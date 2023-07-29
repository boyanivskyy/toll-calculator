package main

import (
	"encoding/json"
	"time"

	"github.com/boyanivskyy/toll-calculator/aggregator/client"
	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer          *kafka.Consumer
	isRunning         bool
	calculatorService CalculatorServicer
	aggregationClient *client.Client
}

func NewKafkaConsumer(topic string, service CalculatorServicer, client *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:          c,
		calculatorService: service,
		aggregationClient: client,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport(consumer) started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	// A signal handler or similar could be used to set this to false to break the loop.
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}

		data := types.OBUData{}
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Error("JSON serialization failed", err)
			continue
		}
		distance, err := c.calculatorService.CalculateDistance(data)
		if err != nil {
			logrus.Error("calculation error", err)
			continue
		}
		reqBody := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}
		if err := c.aggregationClient.AggregateInvoice(reqBody); err != nil {
			logrus.Errorf("aggregate error:", err)
			continue
		}

	}
}
