package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		Secrets  Secrets  `yaml:"secrets"`
		Logger   Logger   `yaml:"logger"`
		Server   Server   `yaml:"server"`
		Postgres Postgres `yaml:"postgres"`
	}

	Secrets struct {
		JwtSecret string `yaml:"jwtSecret"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Env     string `yaml:"env"`
	}

	Logger struct {
		Level string `yaml:"level"`
	}

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		SslMode  string `yaml:"sslmode"`
	}
)

func Init(path string) (*Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
