package command

import (
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var RootCmd = &cmds.Command{Subcommands: map[string]*cmds.Command{}}

func SetHelpText(ht cmds.HelpText) {
	RootCmd.Helptext = ht
}
