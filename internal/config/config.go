package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env"`
	HttpServer HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address"`
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		panic("Error reading config: " + err.Error())
	}

	return &config
}
