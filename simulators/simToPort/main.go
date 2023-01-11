package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"math/rand"
)

type GameEvent struct {
	EventID   string      `json:"eventid"`
	EventType string      `json:"eventtype"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	counter := 0
	for {
		counter++
		event := GameEvent{
			EventID:   fmt.Sprintf("eid%010d", counter),
			Timestamp: time.Now().Unix(),
			EventType: "Spawn",
			Data: map[string]interface{}{
				"playerid": fmt.Sprintf("eid%010d", rand.Intn(100)),
				"x_coord": rand.Intn(90),
				"y_coord": rand.Intn(90),
				"z_coord": rand.Intn(90),
			},
		}

		jsonStr, err := json.Marshal(event)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(string(jsonStr))

		conn.Write(jsonStr)
		time.Sleep(time.Second * 5)
	}
}
