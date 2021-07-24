package main

import (
	"bufio"
	"log"
	"net/url"
	"os"

	"talk/talk_client_speech/alarm"
	"talk/talk_client_speech/job"
	"talk/talk_client_speech/tts"

	"github.com/gorilla/websocket"
)

const WITH_SPAM = true
const INTERVAL_SECONDS_ALARM = 4

func main() {
	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/api/v1/msg"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)

	//TTS
	speaker := tts.New()

	//Creating Worker
	newMsgNotification := func() {
		speaker.Speak("Hay mensajes nuevos")
	}

	// Alarm
	alarm := alarm.New(
		INTERVAL_SECONDS_ALARM,
		newMsgNotification,
	)
	if WITH_SPAM {
		newMsgNotification = alarm.Start
	}
	wm := job.New(
		newMsgNotification,
		func() {
			speaker.Speak("Usted no posee mensajes")
		},
	)

	defer c.Close()

	// Receive messages
	go func() {
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Message received: %s", message)
			job := func() {
				speaker.Speak(string(message))
			}
			wm.AddJob(job)
		}
	}()

	// Read from stdin and send through websocket
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result := scanner.Text()
		if result == "read" {
			if WITH_SPAM {
				alarm.Stop()
				speaker.Stop()
			}
			wm.ExecuteJob()
		} else {
			_ = c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
			log.Printf("Message sent: %s", scanner.Text())
		}
	}
}
