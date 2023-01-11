package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type payload struct {
    EventType string      `json:"event_type"`
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`
}

func main() {
    conn, _ := net.Dial("tcp", "localhost:8080")
    defer conn.Close()
    for {
        event := payload{
            EventType: "example_event",
            Data: map[string]string{
                "example_key": "example_value",
            },
            Timestamp: time.Now().Unix(),
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
