package main

import (
	"github.com/bdengine/ipfs-cmd/cmd"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

func main() {
	err := cmd.Register(map[string]*cmds.Command{"hello": test})
	if err != nil {
		panic(err)
	}
	env, err := ConstructEnv()
	if err != nil {
		panic(err)
	}
	err = cmd.MainRet(env)
	if err != nil {
		panic(err)
	}
}
