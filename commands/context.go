package commands

import (
	"context"
	"fmt"
	"github.com/bdengine/ipfs-cmd/conf"
	cmds "github.com/ipfs/go-ipfs-cmds"
	logging "github.com/ipfs/go-log"
)

const configFileName = "config"

var log = logging.Logger("commands")

type Context struct {
	Ctx        context.Context
	ConfigRoot string
	Config     *conf.Config
	ShutDown   func()
	//server core.Server
}

func ConstructContext(path string) (*Context, error) {
	cfg, err := conf.ReadConfig(path + configFileName)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		Ctx:        ctx,
		ConfigRoot: path,
		Config:     cfg,
		ShutDown:   cancel,
	}, nil
}

func ShutDown(env cmds.Environment) error {
	c, ok := env.(*Context)
	if !ok {
		return fmt.Errorf("不能转化为commands.Context")
	}
	c.ShutDown()
	return nil
}

func Run(env cmds.Environment, ech chan error) error {
	c, ok := env.(*Context)
	if !ok {
		return fmt.Errorf("不能转化为commands.Context")
	}
	select {
	case <-c.Ctx.Done():
		log.Info("接收到退出信号")
		return nil
	case err := <-ech:
		return err
	}
}
