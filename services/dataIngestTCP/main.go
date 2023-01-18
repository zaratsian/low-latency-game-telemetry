package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	enablePrint = true
)

type GameEvent struct {
	EventID   string      `json:"eventid"`
	EventType string      `json:"eventtype"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ln, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			panic(err)
		}
		defer ln.Close()

		counter := 0

		for {
			conn, _ := ln.Accept()
			go func(conn net.Conn) {
				defer conn.Close()

				start_time := time.Now().UnixMilli()
				fmt.Printf("start_time:%v\n", start_time)

				// loop until an error or EOF is encountered
				for {
					buffer := make([]byte, 1024)
					n, err := conn.Read(buffer)
					if err != nil {
						break
					}
					jsonData := buffer[:n]

					var ge GameEvent
					if err := json.Unmarshal(jsonData, &ge); err != nil {
						fmt.Println("Error parsing json:", err)
						return
					}
					counter++

					if enablePrint {
						// process payload for low-latency
						fmt.Println("EventType:", ge.EventType)
						fmt.Println("Timestamp:", ge.Timestamp)
						fmt.Println("Data:", ge.Data)
					}

				}

				fmt.Printf("Event Counter: %v\n", counter)
				fmt.Printf("Runtime:    %v milliseconds\n", time.Now().UnixMilli()-start_time)

			}(conn)

		}
	}()
	wg.Wait()
}
