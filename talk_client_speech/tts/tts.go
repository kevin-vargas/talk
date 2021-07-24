package tts

import (
	"log"
	"os/exec"
	"sync"
)

func New() *TTS {
	auxChan := make(chan bool)
	return &TTS{
		cancelChan: &auxChan,
	}
}

type TTS struct {
	sync.Mutex
	cancelChan *chan bool
}

const (
	PID          = "espeak"
	COMMAND_KILL = "killall"
)

var (
	ARGS = []string{"-s", "140", "-v", "es-mx", "-g", "4"}
)

func (tts *TTS) Speak(text string) {
	args := append(ARGS, text)
	cmd := exec.Command(PID, args...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-*tts.cancelChan:
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	case err := <-done:
		if err != nil {
			log.Fatalf("process finished with error = %v", err)
		}
	}
}

func (tts *TTS) Stop() {
	tts.Lock()
	//BroadCast close
	close(*tts.cancelChan)
	*tts.cancelChan = make(chan bool)
	tts.Unlock()
}
