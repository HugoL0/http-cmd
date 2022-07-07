package environment

import "context"

type Environment interface {
	GetPath() string
	GetPort() string
	GetLogLevel() map[string]string
	GetBin() string
	GetBinEnv() []string
	Daemon() error
}

type Env struct {
	Ctx    context.Context
	Cancel func()
	Environment
}

func ConstructEnv(e Environment) (*Env, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Env{
		Ctx:         ctx,
		Cancel:      cancel,
		Environment: e,
	}, nil
}
