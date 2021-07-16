package main

import (
	"bufio"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/api/v1/msg"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	defer c.Close()

	// Receive messages
	go func() {
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Message received: %s", message)
		}
	}()

	// Read from stdin and send through websocket
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_ = c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		log.Printf("Message sent: %s", scanner.Text())
	}
}
