package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	wxDir = "/Users/hyx/work/gowork/src/gospacex/httpx/examples/wx"
)

// Message represents a chat message
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	User    string `json:"user"`
	Time    string `json:"time,omitempty"`
}

var (
	content = flag.String("content", "Hello from client", "Message content")
	user    = flag.String("user", "Client", "User name")
)

func loadConfig() int {
	viper.Reset()
	viper.SetConfigFile(wxDir + "/config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Failed to read config: %v", err))
	}
	return viper.GetInt("port")
}

func main() {
	flag.Parse()

	port := loadConfig()
	serverURL := fmt.Sprintf("ws://localhost:%d/ws", port)

	log.Printf("Connecting to %s", serverURL)
	log.Printf("Message: content=%q user=%q", *content, *user)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected successfully")

	msg := Message{
		Type:    "message",
		Content: *content,
		User:    *user,
	}

	log.Printf("Sending message: %+v", msg)
	if err := conn.WriteJSON(msg); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent, waiting for broadcasts...")

	for {
		var received Message
		err := conn.ReadJSON(&received)
		if err != nil {
			log.Printf("Connection closed: %v", err)
			break
		}
		log.Printf("Received broadcast: %+v", received)
	}

	log.Println("Client exiting")
}