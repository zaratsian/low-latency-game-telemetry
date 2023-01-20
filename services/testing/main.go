package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	eventCount  int
	//startTime   time.Time
	interval    = 5 * time.Second
	numRoutines = 5000
)

type GameEvent struct {
	Eventid  string
	Playerid string
	Score    float64
}

func main() {

	//startTime = time.Now()

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

func processEvent(ge GameEvent) {
	// Process Event
	eventCount++
}

func (ge GameEvent) validate() bool {
	// Validate Payload
	return true
}
