package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gospacex/httpx"
	_ "github.com/gospacex/httpx/adapter/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message represents a chat message
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	User    string `json:"user"`
	Time    string `json:"time"`
}

// ChatRoom manages connected clients
type ChatRoom struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (cr *ChatRoom) run() {
	for {
		select {
		case conn := <-cr.register:
			cr.mu.Lock()
			cr.clients[conn] = true
			cr.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(cr.clients))

		case conn := <-cr.unregister:
			cr.mu.Lock()
			if _, ok := cr.clients[conn]; ok {
				delete(cr.clients, conn)
				conn.Close()
			}
			cr.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(cr.clients))

		case msg := <-cr.broadcast:
			cr.mu.Lock()
			for conn := range cr.clients {
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					conn.Close()
					delete(cr.clients, conn)
				}
			}
			cr.mu.Unlock()
		}
	}
}

// syncWriter wraps a writer and syncs after every write
type syncWriter struct {
	w *os.File
}

func (sw *syncWriter) Write(p []byte) (int, error) {
	n, err := sw.w.Write(p)
	sw.w.Sync()
	return n, err
}

func wsHandler(room *ChatRoom) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		req := hc.Request().(*http.Request)
		resp := hc.Response().(http.ResponseWriter)
		log.Printf("[WS] Connection from %s %s", req.RemoteAddr, req.URL.Path)
		conn, err := upgrader.Upgrade(resp, req, nil)
		if err != nil {
			log.Printf("[WS] Upgrade failed: %v", err)
			return
		}
		log.Printf("[WS] Connected: %s", req.RemoteAddr)

		room.register <- conn

		go func() {
			defer func() {
				room.unregister <- conn
				conn.Close()
				log.Printf("[WS] Goroutine ended for %s", req.RemoteAddr)
			}()

			for {
				log.Printf("[WS] Waiting for message from %s...", req.RemoteAddr)
				var msg Message
				err := conn.ReadJSON(&msg)
				if err != nil {
					log.Printf("[WS] ReadJSON error from %s: %v (type=%T)", req.RemoteAddr, err, err)
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("[WS] Unexpected close from %s", req.RemoteAddr)
					}
					return
				}

				log.Printf("[WS] RECEIVED from %s: %+v", req.RemoteAddr, msg)

				msg.Time = time.Now().Format(time.RFC3339)
				if msg.User == "" {
					msg.User = "Anonymous"
				}

				log.Printf("[WS] Broadcasting to %d clients: %s", len(room.clients), msg.Content)
				room.broadcast <- msg
			}
		}()
	}
}

func main() {
	// Create log file for test capture
	logFile, err := os.Create("/tmp/wx_server_stdout.log")
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}

	// Use both file and stdout for logging with sync after each write
	log.SetOutput(io.MultiWriter(&syncWriter{w: logFile}, os.Stdout))

	log.Printf("Server starting...")

	room := newChatRoom()
	go room.run()

	app, err := httpx.New("/Users/hyx/work/gowork/src/gospacex/httpx/examples/wx/config.yaml")
	if err != nil {
		panic(err)
	}

	// WebSocket endpoint
	app.WS("/ws", wsHandler(room))

	// HTTP endpoints for info
	app.GET("/", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]interface{}{
			"service": "httpx WebSocket Demo",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"WebSocket": "ws://localhost:8081/ws",
				"HTTP":      "http://localhost:8081/",
			},
			"instructions": map[string]string{
				"connect":    "Connect to ws://localhost:8081/ws using WebSocket client",
				"send":       `Send JSON: {"type":"message","content":"Hello","user":"YourName"}`,
				"message":    "Server broadcasts to all connected clients",
			},
		})
	})

	app.GET("/status", func(ctx context.Context, hc httpx.HandlerContext) {
		room.mu.Lock()
		count := len(room.clients)
		room.mu.Unlock()
		hc.AbortJSON(200, map[string]interface{}{
			"clients": count,
			"uptime":  "running",
		})
	})

	fmt.Print(`
╔════════════════════════════════════════════════════════╗
║     httpx WebSocket Demo                               ║
║     WebSocket: ws://localhost:8081/ws                  ║
║     HTTP:      http://localhost:8081/                  ║
╚════════════════════════════════════════════════════════╝
Test with JavaScript:
  const ws = new WebSocket('ws://localhost:8081/ws');
  ws.onopen = () => ws.send(JSON.stringify({type:'message',content:'Hello',user:'Test'}));
  ws.onmessage = (e) => console.log(JSON.parse(e.data));
`)
	if err := app.Run(); err != nil {
		panic(err)
	}
}