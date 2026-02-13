package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config with tags from cleanenv library
type ServerConfig struct {
	Host           string        `env:"API_HOST" env-required:"true"`
	Port           string        `env:"API_PORT" env-required:"true"`
	MaxHeaderBytes int           `yaml:"max_header_bytes" env-default:"1048576"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env-default:"10s"`
}

// Database config from env and config.yaml
type DatabaseConfig struct {
	Username string `env:"DATABASE_USERNAME" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

// Logger config from config.yaml
type LoggerConfig struct {
	Level  string `yaml:"level" env-default:"info"`
	Format string `yaml:"format" env-default:"json"`
}

// Dataclass with all configs
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

// Load config from config/config.yaml
func LoadConfig(path string) (Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return Config{}, err
	}

	switch config.Logger.Level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel:
	default:
		return Config{}, fmt.Errorf(`selected incorrect logger level: "%s"`, config.Logger.Level)
	}

	switch config.Logger.Format {
	case TextFormat, JsonFormat:
	default:
		return Config{}, fmt.Errorf(`selected incorrect logger format: "%s"`, config.Logger.Format)
	}

	switch config.Database.SSLMode {
	case SSLDisable, SSLRequire:
	default:
		return Config{}, fmt.Errorf(`selected incorrect ssl mode: "%s"`, config.Database.SSLMode)
	}

	return config, nil
}
