package main

import (
    "log"
    "net/http"
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

type Packet struct {
    ID        string  `json:"id"`
    Timestamp float64 `json:"timestamp"`  // Changed to float64 for JS timestamp
    Payload   string  `json:"payload"`
    Size      int     `json:"size"`
    Latency   float64 `json:"latency,omitempty"`
}

type Client struct {
    conn *websocket.Conn
    mu   sync.Mutex
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (c *Client) writeJSON(v interface{}) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.conn.WriteJSON(v)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Upgrade error: %v", err)
        return
    }

    client := &Client{conn: conn}
    defer client.conn.Close()

    for {
        var packet Packet
        err := conn.ReadJSON(&packet)
        if err != nil {
            if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
                log.Printf("Read error: %v", err)
            }
            break
        }

        now := float64(time.Now().UnixNano()) / 1e6  // Convert to JS timestamp (milliseconds)
        packet.Latency = now - packet.Timestamp
        packet.Timestamp = now

        err = client.writeJSON(packet)
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
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
