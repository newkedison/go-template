// +build !windows
package service

import (
	"os"
	"os/signal"
	"runtime"
)

type systemNoService struct {
}

func (s systemNoService) Init(option ServiceOption) error {
	return nil
}

func (s systemNoService) Interactive() bool {
	return true
}

func (s systemNoService) Platform() string {
	return runtime.GOOS
}

func (s systemNoService) Run(f func(exit <-chan struct{})) error {
	exit := make(chan struct{})
	go f(exit)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	close(exit)
	return nil
}

func init() {
	_system = systemNoService{}
}
