package register

import (
	"fmt"
	"github.com/bdengine/ipfs-cmd/cmd"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var CmdMap = map[string]*cmds.Command{
	"daemon":   cmd.DaemonCmd,
	"shutdown": cmd.ShutdownCmd,
}

func Register(cmdMap map[string]*cmds.Command) error {
	for s, command := range cmdMap {
		if _, ok := CmdMap[s]; ok {
			return fmt.Errorf("cmd %s repeat", s)
		}
		CmdMap[s] = command
	}
	return nil
}
