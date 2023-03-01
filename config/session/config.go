package session

import (
	"engine/config/config"
	"strings"
)

type driver string

const (
	Memory driver = "memory"
	Redis  driver = "redis"
)

var defaultName = "score-system-session"

type sessionConfig struct {
	Driver string
	Name   string
}

func getConfig() sessionConfig {

	wc := sessionConfig{}
	wc.Driver = string(Memory)
	if config.Config.IsSet("session.driver") {
		wc.Driver = strings.ToLower(config.Config.GetString("session.driver"))
	}

	wc.Name = defaultName
	if config.Config.IsSet("session.name") {
		wc.Name = strings.ToLower(config.Config.GetString("session.name"))
	}

	return wc
}
