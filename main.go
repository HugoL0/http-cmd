package main

import (
	"github.com/bdengine/ipfs-cmd/cmd"
	logging "github.com/ipfs/go-log"
	"os"
)

var log = logging.Logger("main")

func main() {
	err := cmd.MainRet()
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}
}
