package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for testing
    },
}

type Packet struct {
    Timestamp int64  `json:"timestamp"`
    Payload   string `json:"payload"`
    Size      int    `json:"size"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Upgrade error: %v", err)
        return
    }
    defer conn.Close()

    for {
        var packet Packet
        err := conn.ReadJSON(&packet)
        if err != nil {
            log.Printf("Read error: %v", err)
            break
        }

        // Echo back the packet
        packet.Timestamp = time.Now().UnixNano()
        err = conn.WriteJSON(packet)
        if err != nil {
            log.Printf("Write error: %v", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)
    http.Handle("/", http.FileServer(http.Dir("static")))
    
    log.Printf("Server starting on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
