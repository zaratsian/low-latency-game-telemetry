package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	// UDP Listener
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Continuosly read packets
	for {
		handleJSON(conn)
	}
}

func handleJSON(conn *net.UDPConn) {
	// Buffer packets
	buffer := make([]byte, 2048)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	// Process JSON data
	var data map[string]interface{}
	err = json.Unmarshal(buffer[:n], &data)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(data)
}
