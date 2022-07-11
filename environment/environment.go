package environment

import "context"

type Environment interface {
	// GetPath 获取工作路径
	GetPath() string
	// GetPort 获取监听端口
	GetPort() string
	// GetLogLevel 日志级别定义
	GetLogLevel() map[string]string
	// GetBin 获取命令路径（主要用于daemon  --background=true（默认）命令时，找到命令）
	GetBin() string
	// GetBinEnv 命令环境变量
	GetBinEnv() []string
	// Daemon daemon命令执行时，用户自定义的daemon命令
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
