package config

import (
	"net/http"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Server config with tags from cleanenv library
type ServerConfig struct {
	Host           string        `yaml:"-"`
	Port           string        `yaml:"-"`
	Handler        http.Handler  `yaml:"-"`
	MaxHeaderBytes int           `yaml:"max_header_bytes" env-default:"1048576"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env-default:"10s"`
}

// Load config from config/config.yaml
func LoadConfig(path string) (ServerConfig, error) {
	var config ServerConfig
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return ServerConfig{}, err
	}

	return config, nil
}
