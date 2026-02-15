package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config with tags from cleanenv library
type ServerConfig struct {
	Host                    string        `env:"API_HOST" env-required:"true"`
	Port                    string        `env:"API_PORT" env-required:"true"`
	MaxHeaderBytes          int           `yaml:"max_header_bytes" env-default:"1048576"`
	ReadTimeout             time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout            time.Duration `yaml:"write_timeout" env-default:"10s"`
	TimeForGracefulShutdown time.Duration `yaml:"time_for_graceful_shutdown" env-default:"10s"`
}

// Database config from env and config.yaml
type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `env:"DATABASE_PORT" env-required:"true"`
	Username string `env:"DATABASE_USERNAME" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true" validate:"oneof=disable require"`
}

// Logger config from config.yaml
type LoggerConfig struct {
	Level  string `yaml:"level" env-default:"info" validate:"oneof=debug info warn error"`
	Format string `yaml:"format" env-default:"json" validate:"oneof=text json"`
}

// Parser config from config.yaml
type ParserConfig struct {
	CountOfWorkers int           `yaml:"count_of_workers" env-default:"3" validate:"gte=1,lte=10"`
	ScanFrequency  time.Duration `yaml:"scan_frequency" env-default:"60s"`
	InputDir       string        `yaml:"input_dir" env-default:"./input-data"`
	OutputDir      string        `yaml:"output_dir" env-default:"./output-data"`
	ErrorDir       string        `yaml:"error_dir" env-default:"./error-data"`
}

// Dataclass with all configs
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
	Parser   ParserConfig   `yaml:"parser"`
}

// Load config from config/config.yaml
func LoadConfig(path string) (Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return Config{}, err
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return Config{}, err
	}

	return config, nil
}
