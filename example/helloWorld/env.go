package main

import (
	"fmt"
	"github.com/hugo/http-cmd/environment"
)

type Config struct {
	Port     string
	Path     string
	LogLevel map[string]string
	Bin      string
	BinEnv   []string
}

type MyEnv struct {
	config Config
}

func (e MyEnv) Daemon() error {
	fmt.Printf("api listen at %s \n", e.config.Port)
	return nil
}

func ConstructEnv() (*environment.Env, error) {
	m := &MyEnv{
		config: Config{
			Port:     "127.0.0.1:9999",
			Path:     "D:/code/go/src/github.com/hugo/http-cmd/example/helloWorld",
			LogLevel: nil,
			Bin:      "D:/code/go/src/github.com/hugo/http-cmd/example/helloWorld/hello.exe",
			BinEnv:   nil,
		},
	}
	return environment.ConstructEnv(m)
}

func (e MyEnv) GetPath() string {
	return e.config.Path
}

func (e MyEnv) GetPort() string {
	return e.config.Port
}

func (e MyEnv) GetLogLevel() map[string]string {
	return e.config.LogLevel
}

func (e MyEnv) GetBin() string {
	return e.config.Bin
}

func (e MyEnv) GetBinEnv() []string {
	return e.config.BinEnv
}
