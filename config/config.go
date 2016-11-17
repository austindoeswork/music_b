package config

import (
	"encoding/json"
	"fmt"
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

func (c *Config) PPrint() {
	fmt.Printf("CONFIG:\nServerPath: %s\nMusicDir: %s\nStaticDir: %s\n",
		c.ServerPath, c.MusicDir, c.StaticDir)
}
