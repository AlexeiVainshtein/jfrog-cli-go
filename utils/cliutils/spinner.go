package cliutils

import (
	"time"
	"fmt"
	"sync"
	"os"
)

type Spinner struct {
	Delay    time.Duration
	lock     *sync.RWMutex
	prefix   string
	active   bool
	stopChan chan struct{}
}

func NewSpinner(prefix string, d time.Duration) *Spinner {
	return &Spinner{
		Delay:    d,
		lock:     &sync.RWMutex{},
		prefix:   prefix,
		active:   false,
		stopChan: make(chan struct{}, 1),
	}
}

func (s *Spinner) Start() {
	if s.active || os.Getenv("JFROG_CLI_LOG_LEVEL") == "ERROR" {
		return
	}
	s.active = true
	fmt.Print(s.prefix)
	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.lock.Lock()
				fmt.Print(".")
				delay := s.Delay
				s.lock.Unlock()
				time.Sleep(delay)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if os.Getenv("JFROG_CLI_LOG_LEVEL") == "ERROR" {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.active {
		s.active = false
		fmt.Print(" Done.\n")
		s.stopChan <- struct{}{}
	}
}