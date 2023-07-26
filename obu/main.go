package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/gorilla/websocket"
)

const (
	sendInverval = time.Second * 5
	wsEndpoint   = "ws://127.0.0.1:30000/ws"
)

func genCords() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()

	return n + f
}
func genLatLong() (float64, float64) {
	return genCords(), genCords()
}

func generateOBUIds(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func main() {
	obuIds := generateOBUIds(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIds[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInverval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
