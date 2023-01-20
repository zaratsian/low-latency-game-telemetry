package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
    udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        panic(err)
    }

    udpConn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        panic(err)
    }
    defer udpConn.Close()

	for {
		processData(udpConn)
	}
}

func processData(conn *net.UDPConn) {
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