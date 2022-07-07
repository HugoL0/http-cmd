package cmd

import (
	"fmt"
	"github.com/bdengine/ipfs-cmd/environment"
	"github.com/bdengine/ipfs-cmd/lock"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var ShutdownCmd = &cmds.Command{
	Run: func(req *cmds.Request, emit cmds.ResponseEmitter, env cmds.Environment) error {
		ctx := env.(environment.Environment)
		locked, err := lock.CheckLocked(ctx.GetPath())
		if err != nil {
			return err
		}
		if !locked {
			return fmt.Errorf("守护进程没有运行")
		}
		return ShutDown(env)
	},
	NoLocal: true,
}

func ShutDown(env cmds.Environment) error {
	c, ok := env.(*environment.Env)
	if !ok {
		return fmt.Errorf("不能转化为commands.Context")
	}
	c.Cancel()
	return nil
}
