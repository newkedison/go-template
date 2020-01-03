package service

import (
	"errors"
	"io/ioutil"
	"strings"
	"os"
	"os/signal"

	"TEMPLATE/utils/logrus/hooks"

	_service "github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

var config ServiceOption

type program struct {
	exit chan struct{}
	f    func(exit <-chan struct{})
}

type systemWindows struct {
}

func (s systemWindows) Init(option ServiceOption) error {
	option.Name = strings.TrimSpace(option.Name)
	if len(option.Name) == 0 {
		return errors.New("Name must not be empty")
	}
	if strings.Index(option.Name, " ") >= 0 {
		return errors.New("Name must not contain space")
	}
	config = option
	return nil
}

func (s systemWindows) Interactive() bool {
	return _service.Interactive()
}

func (s systemWindows) Platform() string {
	return _service.Platform()
}

func (sys systemWindows) Run(f func(exit <-chan struct{})) error {
	if config.Name == "" {
		exit := make(chan struct{})
		go f(exit)
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)
		<-signalCh
		close(exit)
		return nil
	}
	p := &program{
		exit: make(chan struct{}),
		f:    f,
	}
	cfg := &_service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}
	s, err := _service.New(p, cfg)
	if err != nil {
		return err
	}
	if !_service.Interactive() {
		logrus.SetOutput(ioutil.Discard)
		formatter := logrus.StandardLogger().Formatter
		switch v := formatter.(type) {
		case *logrus.TextFormatter:
			v.DisableTimestamp = true
		case *logrus.JSONFormatter:
			v.DisableTimestamp = true
		default:
			// do nothing
		}
		if logger, err := s.Logger(nil); err == nil {
			logrus.AddHook(hooks.NewServiceHook(logger))
		}
	}
	return s.Run()
}

func (p *program) Start(s _service.Service) error {
	go p.f(p.exit)
	return nil
}

func (p *program) Stop(s _service.Service) error {
	close(p.exit)
	return nil
}

func init() {
	_system = systemWindows{}
}
