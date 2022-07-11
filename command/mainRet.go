package command

import (
	"context"
	"fmt"
	"github.com/hugo/http-cmd/environment"
	"github.com/hugo/http-cmd/lock"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs-cmds/cli"
	"github.com/ipfs/go-ipfs-cmds/http"
	"os"
)

func MainRet(env *environment.Env) error {
	err := constructCmd()
	if err != nil {
		return err
	}
	// parse the command path, arguments and options from the command line
	req, err := cli.Parse(context.TODO(), os.Args[1:], os.Stdin, RootCmd)
	if err != nil {
		return err
	}
	//req.Options["encoding"] = cmds.Text
	// create an emitter
	res, err := cli.NewResponseEmitter(os.Stdout, os.Stderr, req)
	if err != nil {
		return err
	}
	executor, err := makeExecutor(req, env)
	if err != nil {
		return err
	}
	// send request to server
	err = executor.Execute(req, res, env)
	if err != nil {
		return err
	}
	return nil
}

func constructCmd() error {
	for s, command := range cmdMap {
		// 检查是否重复
		if _, ok := RootCmd.Subcommands[s]; ok {
			return fmt.Errorf("重复的命令%v", s)
		}
		RootCmd.Subcommands[s] = command
	}
	return nil
}

func makeExecutor(req *cmds.Request, env environment.Environment) (cmds.Executor, error) {
	rok, err := checkRemote(req.Command, env.GetPath())
	if err != nil {
		return nil, err
	}
	var executor cmds.Executor
	if rok {
		// create http rpc client
		executor = http.NewClient(env.GetPort())
	} else {
		executor = cmds.NewExecutor(RootCmd)
	}
	return executor, nil
}

func checkRemote(c *cmds.Command, configRoot string) (bool, error) {
	if c.NoRemote && c.NoLocal {
		return false, fmt.Errorf("命令设置出错,noremote && nolocal")
	}
	locked, err := lock.CheckLocked(configRoot)
	if err != nil {
		return false, err
	}
	if c.NoLocal {
		if !locked {
			return false, fmt.Errorf("守护进程未运行，无法执行该命令")
		}
		return true, nil
	}
	if c.NoRemote {
		if locked {
			return false, fmt.Errorf("不能多进程的执行local命令")
		}
		return false, nil
	}
	return locked, nil
}
