package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	CONFIG_PATH = "./config/config.yaml"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	Storage     string `yaml:"storage" env-default:"memory"`
	AliasLength int    `yaml:"alias_length" env-default:"10"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

func MustLoad() *Config {
	configPath := CONFIG_PATH
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist: %s", configPath))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("cannot read config: %s", err))
	}

	return &cfg
}
