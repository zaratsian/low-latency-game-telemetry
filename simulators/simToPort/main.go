package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type GameEvent struct {
	EventID	  string 		`json:"eventid"`
    EventType string 		`json:"eventtype"`
	Timestamp int64			`json:"timestamp"`
    Data      interface{} 	`json:"data"`
}

func main() {
    conn, _ := net.Dial("tcp", "localhost:8080")
    defer conn.Close()
    for {
        event := GameEvent{
			EventID: "idxyz123",
            Timestamp: time.Now().Unix(),
			EventType: "example_event",
			Data: map[string]string{
                "example_key": "example_value",
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
