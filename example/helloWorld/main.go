package main

import (
	"github.com/hugo/http-cmd/command"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

func main() {
	err := command.Register(map[string]*cmds.Command{"hello": test})
	if err != nil {
		panic(err)
	}
	env, err := ConstructEnv()
	if err != nil {
		panic(err)
	}
	err = command.MainRet(env)
	if err != nil {
		panic(err)
	}
}
