package hooks

import (
	"fmt"
	"os"

	_service "github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

// ServiceHook to send logs via github.com/kardianos/service.
type ServiceHook struct {
	l _service.Logger
}

func NewServiceHook(l _service.Logger) *ServiceHook {
	return &ServiceHook{l}
}

func (hook *ServiceHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return hook.l.Error(line)
	case logrus.WarnLevel:
		return hook.l.Warning(line)
	case logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel:
		return hook.l.Info(line)
	default:
		return nil
	}
}

func (hook *ServiceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
