package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type GRPCConfig struct {
	Port    int8          `yaml:"port" env-default:"15900"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type Config struct {
	TokenLifetime time.Duration `yaml:"token_lifetime" env-required:"true"`
	StoragePath   string        `yaml:"storage_path" env-required:"true"`
	GRPC          GRPCConfig    `yaml:"grpc" env-required:"true"`
}

/*
* load config from path
* check if path is valid (file exists)
* read config by cleanenv.ReadConfig(path string, cfg interface{})
* error handling
 */

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	return nil, nil
	cleanenv.ReadConfig()
}

// get path from env or flag or default
func getConfigPath() (string, error) {
	return "", nil
}
