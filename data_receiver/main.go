package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/gorilla/websocket"
)

func main() {
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
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
	fmt.Println("New OBU client(data_receiver) connected!")
	for {
		data := types.OBUData{}
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read ws data error:", err)
			continue
		}
		fmt.Printf("received OBU data from [%d]:: <lat %.2f, long %.2f> \n", data.OBUID, data.Lat, data.Long)

		dr.msgch <- data
	}
}
