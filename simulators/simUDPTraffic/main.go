package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	concurrency = 200
	timeDelay   = 1 // milliseconds
)

type GameEvent struct {
	Eventid  string 	`json:"eventid"`
	Datetime string 	`json:"datetime"`
	Playerid string		`json:"playerid"`
	Event string        `json:"event"`
	Score float64       `json:"score"`
}

var wg sync.WaitGroup

func main() {

	simulations := flag.Int64("simulations", 1000000, "Number of events to simulate")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	addr, _ := net.ResolveUDPAddr("udp", "localhost:8080")
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	wg.Add(1)
	go sendTraffic(conn, *simulations)
	wg.Wait()
}

func sendTraffic(conn *net.UDPConn, sims int64) {
	defer wg.Done()

	numRoutines := concurrency
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for simNum := 0; simNum <= int(sims/concurrency); simNum++ {
				ge := generateGameEvent()
				jsonEvent, _ := json.Marshal(ge)
				conn.Write(jsonEvent)
				time.Sleep(time.Duration(rand.Intn(timeDelay)) * time.Millisecond)
			}
		}()
	}
}

func generateGameEvent() GameEvent {
	return GameEvent{
		Eventid: fmt.Sprintf("event_%10d", rand.Intn(1000000000)), 
		Playerid: fmt.Sprintf("player_%06d", rand.Intn(100000)),
		Datetime: fmt.Sprintf("2022-%02d-%02d", rand.Intn(12), rand.Intn(28)),
		Event: "action",
		Score: rand.Float64() * 100}
}
