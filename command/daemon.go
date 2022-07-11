package command

import (
	"fmt"
	"github.com/hugo/http-cmd/environment"
	"github.com/hugo/http-cmd/lock"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs-cmds/http"
	logging "github.com/ipfs/go-log"
	nethttp "net/http"
	"os/exec"
)

var log = logging.Logger("cmd")

const (
	backgroundOpt = "background"
	restartOpt    = "reStart"
)

var DaemonCmd = &cmds.Command{
	Options: []cmds.Option{
		cmds.BoolOption(backgroundOpt, "是否后台运行").WithDefault(true),
		cmds.BoolOption(restartOpt, "是否强制重启").WithDefault(false),
	},
	Arguments: nil,
	PreRun:    nil,
	Run: func(req *cmds.Request, emit cmds.ResponseEmitter, env cmds.Environment) error {
		c := env.(*environment.Env)
		locked, err := lock.CheckLocked(c.GetPath())
		if err != nil {
			return err
		}
		if locked {
			return fmt.Errorf("守护进程正在运行中")
		}
		for s1, s2 := range c.GetLogLevel() {
			err := logging.SetLogLevel(s1, s2)
			if err != nil {
				log.Errorf("error setting %s to level %s,err:%s", s1, s2, err)
			}
		}
		background := req.Options[backgroundOpt].(bool)
		if background {
			// 启动子命令
			e := &exec.Cmd{
				Path: c.GetBin(),
				Args: []string{c.GetBin(), "daemon", "--background=false"},
				Env:  c.GetBinEnv(),
			}
			err := e.Start()
			if err != nil {
				return err
			}
			fmt.Println("启动 守护进程成功")
			return nil
		}
		log.Info("守护进程 启动中")
		f, err := lock.TryLockDaemon("")
		if err != nil {
			return err
		}
		defer f.Close()

		ech := make(chan error, 5)
		go func() {
			log.Infof("启动http服务,端口%s", c.GetPort())
			h := http.NewHandler(env, RootCmd, http.NewServerConfig())
			// create http rpc server
			err = nethttp.ListenAndServe(c.GetPort(), h)
			if err != nil {
				ech <- err
				return
			}
		}()

		select {
		case <-c.Ctx.Done():
			log.Info("接收到退出信号")
			return nil
		case err := <-ech:
			return err
		default:
			go func() {
				err = c.Daemon()
				if err != nil {
					ech <- err
				}
			}()
		}

		err = Run(c, ech)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	},
	NoRemote: true,
}

func Run(env *environment.Env, ech chan error) error {
	select {
	case <-env.Ctx.Done():
		log.Info("接收到退出信号")
		return nil
	case err := <-ech:
		return err
	}
}
