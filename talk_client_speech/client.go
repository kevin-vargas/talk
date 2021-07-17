package main

import (
	"bufio"
	"log"
	"net/url"
	"os"

	"talk/talk_client_speech/job"
	"talk/talk_client_speech/tts"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/api/v1/msg"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	//Creating Worker
	wm := job.New()

	defer c.Close()

	// Receive messages
	go func() {
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Message received: %s", message)
			job := func() {
				tts.Speak(string(message))
			}
			wm.AddJob(job)
		}
	}()

	// Read from stdin and send through websocket
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result := scanner.Text()
		if result == "read" {
			wm.ExecuteJob()
		} else {
			_ = c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
			log.Printf("Message sent: %s", scanner.Text())
		}
	}
}
