package command

import (
	"fmt"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

const (
	cmdDaemon   = "daemon"
	cmdShutdown = "shutdown"
)

var cmdMap = map[string]*cmds.Command{
	cmdDaemon:   DaemonCmd,
	cmdShutdown: ShutdownCmd,
}

func Register(cm map[string]*cmds.Command) error {
	for s, command := range cm {
		if _, ok := cmdMap[s]; ok {
			return fmt.Errorf("命令 %s 重复", s)
		}
		cmdMap[s] = command
	}
	return nil
}

func Unregister(cmdName string) error {
	if _, ok := cmdMap[cmdName]; ok {
		delete(cmdMap, cmdName)
		return nil
	}
	return fmt.Errorf("命令%s不存在", cmdName)
}
