package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Listen for events from Agones
	ln, err := net.Listen("tcp", ":27960")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		// Accept connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go processData(conn)
	}
}

func processData(conn net.Conn) {
	data := make([]byte, 512)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Received data:", string(data))
	// TODO

	conn.Close()
}
