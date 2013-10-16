package multiplexity

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Debug bool

	Server struct {
		Address  string `json:address`
		Nick     string `json:nick`
		User     string `json:user`
		RealName string `json:realName`
	} `json:"server"`

	Client struct {
		ListenAddress string `json:listenAddress`
	} `json:client`
}

func LoadConfig(fileName string) (*Config, error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
