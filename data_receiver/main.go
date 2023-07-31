package main

import (
	"log"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func main() {

	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch    chan types.OBUData
	conn     *websocket.Conn
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLoggingMiddleware(p)

	return &DataReceiver{
		msgch:    make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.producer.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	logrus.Info("New OBU client(data_receiver) connected!")

	for {
		data := types.OBUData{}
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read ws error:", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			log.Println("kafka produce error", err)
		}
	}
}
