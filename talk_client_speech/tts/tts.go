package tts

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	COMMAND = "espeak"
)

var (
	ARGS = []string{"-s", "140", "-v", "es-mx", "-g", "4"}
)

func Speak(text string) {
	args := append(ARGS, text)
	cmd := exec.Command(COMMAND, args...)
	fmt.Println("READING: ", text)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
