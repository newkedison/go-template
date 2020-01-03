package service

var (
	_system  system
)

type system interface {
	Init(option ServiceOption) error
	Interactive() bool
	Platform() string
	Run(func(<-chan struct{})) error
}

type ServiceOption struct {
	Name        string
	DisplayName string
	Description string
}

func Init(option ServiceOption) error {
	return _system.Init(option)
}

func Interactive() bool {
	return _system.Interactive()
}

func Platform() string {
	return _system.Platform()
}

func Run(f func(exit <-chan struct{})) error {
	return _system.Run(f)
}
