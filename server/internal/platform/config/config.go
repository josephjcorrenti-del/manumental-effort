package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Auth    AuthConfig    `yaml:"auth"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type AuthConfig struct {
	JWTSigningKey      string `yaml:"jwt_signing_key"`
	TokenExpiryMinutes int    `yaml:"token_expiry_minutes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
