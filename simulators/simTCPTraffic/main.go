package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	simulations     = 10
	enablePrint     = true
	enableTimeDelay = true
	timeDelay       = 1000 // millisecond(s)
)

type GameEvent struct {
	EventID   string      `json:"eventid"`
	EventType string      `json:"eventtype"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func main() {

	go func() {

		conn, _ := net.Dial("tcp", "localhost:8080")
		defer conn.Close()
		for counter := 0; counter <= simulations; counter++ {
			event := GameEvent{
				EventID:   fmt.Sprintf("eid%010d", counter),
				Timestamp: time.Now().Unix(),
				EventType: "Spawn",
				Data: map[string]interface{}{
					"playerid": fmt.Sprintf("eid%010d", rand.Intn(100)),
					"x_coord":  rand.Intn(90),
					"y_coord":  rand.Intn(90),
					"z_coord":  rand.Intn(90),
				},
			}

			jsonStr, err := json.Marshal(event)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			if enablePrint {
				fmt.Println(string(jsonStr))
			}

			conn.Write(jsonStr)

			if enableTimeDelay {
				time.Sleep(time.Millisecond * timeDelay)
			}
		}

	}()

}
