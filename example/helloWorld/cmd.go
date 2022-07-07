package main

import (
	"fmt"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var test = &cmds.Command{
	Options: []cmds.Option{
		cmds.StringsOption("name", "n", "your name"),
	},
	Arguments: []cmds.Argument{
		{
			Name:          "you",
			Type:          0,
			Required:      false,
			Variadic:      false,
			SupportsStdin: false,
			Recursive:     false,
			Description:   "",
		},
	},
	Run: func(req *cmds.Request, emit cmds.ResponseEmitter, env cmds.Environment) error {
		name, _ := req.Options["name"].(string)
		res := "你好,陌生人"
		if name != "" {
			res = fmt.Sprintf("你好,%s", name)

		}
		return emit.Emit(res)
	},
}
