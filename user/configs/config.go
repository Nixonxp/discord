package config

import (
	"github.com/Nixonxp/discord/user/pkg/configurator"
	"sync"
)

type Config struct {
	Application ApplicationConfig `yaml:"application"`
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) Configurable() {}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = NewConfig()
		configurator.ReadEnvironment(instance)
	})

	return instance
}
