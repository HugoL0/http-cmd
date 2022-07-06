package cmd

import (
	"fmt"
	"github.com/bdengine/ipfs-cmd/commands"
	"github.com/bdengine/ipfs-cmd/lock"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var ShutdownCmd = &cmds.Command{
	Run: func(req *cmds.Request, emit cmds.ResponseEmitter, env cmds.Environment) error {
		ctx := env.(*commands.Context)
		locked, err := lock.CheckLocked(ctx.ConfigRoot)
		if err != nil {
			return err
		}
		if !locked {
			return fmt.Errorf("守护进程没有运行")
		}
		return commands.ShutDown(env)
	},
	NoLocal: true,
}
