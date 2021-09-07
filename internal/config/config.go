package config

import (
	"genesis/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" `
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIp string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8000"`
	} `yaml:"listen"`
}

var instance *Config
var once *sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read configuration")
		instance = &Config{}
		err := cleanenv.ReadConfig("config.yml", instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
