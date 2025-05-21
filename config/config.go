package config

import "github.com/caarlos0/env"

type Config struct {
	POSTGRES_DB string `env:"POSTGRES_DB" envDefault:"postgres"`
	PORT        int    `env:"PORT" envDefault:"8181"`
}

func New() (*Config, error) {
	config := &Config{}

	return config, env.Parse(config)
}
