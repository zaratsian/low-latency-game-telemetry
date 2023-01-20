package main

import (
	"encoding/json"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	simulations = 1000000
	concurrency = 2500
	timeDelay   = 1 // milliseconds
)

type GameEvent struct {
	Eventid  string
	Playerid string
	Score    float64
}

var wg sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())

	addr, _ := net.ResolveUDPAddr("udp", "localhost:8080")
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	wg.Add(1)
	go sendTraffic(conn)
	wg.Wait()
}

func sendTraffic(conn *net.UDPConn) {
	defer wg.Done()

	numRoutines := concurrency
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for simNum := 0; simNum <= int(simulations/concurrency); simNum++ {
				ge := generateGameEvent()
				jsonEvent, _ := json.Marshal(ge)
				conn.Write(jsonEvent)
				time.Sleep(time.Duration(rand.Intn(timeDelay)) * time.Millisecond)
			}
		}()
	}
}

func generateGameEvent() GameEvent {
	return GameEvent{Eventid: "random_string", Playerid: "player123", Score: rand.Float64() * 100}
}
