package main

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println(err)
            return
        }
        defer conn.Close()

        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                log.Println(err)
                break
            }
            log.Printf("Received message: %s", msg)
			
			//placeholder for data msg processing
        }
    })

    log.Println("Listening for game events on ws://localhost:8080/ws")
    log.Fatal(http.ListenAndServe(":8080", nil))
}


