package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	HelperBin string
	HelperEnv []string
	IpfsBin   string
	IpfsEnv   []string
	Port      string

	LogLevel map[string]string
}

func ReadConfig(file string) (*Config, error) {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	res := &Config{}
	err = json.Unmarshal(read, res)
	if err != nil {
		return nil, err
	}
	if res.Port == "" {
		res.Port = "/8888"
	}
	return res, nil
}
