package main

import (
	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
	"github.com/pelletier/go-toml"
)


func LoadConfig(filepath string) (*Config, error) {
	logger.Standard().Infof("Reading config path %s", filepath)

	t, err := toml.LoadFile(filepath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err := t.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
