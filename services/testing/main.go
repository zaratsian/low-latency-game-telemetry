package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/low-latency-game-telemetry/spanner"
	"github.com/low-latency-game-telemetry/utils"
)

var (
	eventCount int
	interval    = 2 * time.Second
	numRoutines = 200
	writeToSpanner = false
)

type GameEvent struct {
	Eventid  string
	Playerid string
	Score    float64
}

func main() {

	addr, _ := net.ResolveUDPAddr("udp", ":8080")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go handleEvents(conn, &wg)
	wg.Wait()
}

func handleEvents(conn *net.UDPConn, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()

	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for {
				var ge GameEvent
				buffer := make([]byte, 1024)
				n, _, _ := conn.ReadFromUDP(buffer)
				json.Unmarshal(buffer[:n], &ge)

				// Validate Payload
				if !ge.validate() {
					fmt.Println("Invalid event received")
					return
				}

				// Process Event
				//go processEvent(ge)

				if writeToSpanner {
					key_string, value_string := utils.FormatStruct(ge)
					err := spanner.SpannerWriteDML(ctx, key_string, value_string)
					if err != nil {
						log.Printf("Error when writing to Spanner. %v\n", err)
					}
				}

				eventCount++
			}
		}()
	}

	ticker := time.NewTicker(interval)
	for range ticker.C {
		fmt.Printf("Events per second: %.3f\n", float64(eventCount)/interval.Seconds())
		eventCount = 0
	}
}

func processEvent(ge GameEvent) bool {
	// Process Event
	return true
}

func (ge GameEvent) validate() bool {
	// Validate Payload
	return true
}