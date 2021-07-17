package alarm

import (
	"sync"
	"time"
)

type Alarm struct {
	sync.Mutex
	notification func()
	interval     int
	stopChan     chan bool
	on           bool
}

func New(interval int, notification func()) *Alarm {
	return &Alarm{
		notification: notification,
		interval:     interval,
		stopChan:     make(chan bool, 1),
		on:           false,
	}
}

func (a *Alarm) Start() {
	a.Lock()
	if a.on {
		return
	}
	a.on = true
	a.Unlock()
L:
	for {
		select {
		case <-a.stopChan:
			break L
		case <-time.After(time.Duration(a.interval) * time.Second):
			a.notification()
		}
	}
}

func (a *Alarm) Stop() {
	a.Lock()
	a.on = false
	a.Unlock()
	a.stopChan <- true
}
