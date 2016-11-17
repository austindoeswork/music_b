package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ServerPath string
	MusicDir   string
	StaticDir  string
}

func Parse(configPath string) (*Config, error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
