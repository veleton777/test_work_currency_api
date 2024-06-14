package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	app          app
	postgres     postgres
	http         http
	fastForexAPI fastForexAPI
}

type app struct {
	Name     string `envconfig:"APP_NAME"`
	LogLevel int8   `envconfig:"LOG_LEVEL"`
}

type http struct {
	Port int32 `envconfig:"HTTP_PORT"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	var cnf Config

	if err := envconfig.Process("", &cnf.app); err != nil {
		return Config{}, errors.Wrap(err, "parse app env")
	}

	if err := envconfig.Process("", &cnf.postgres); err != nil {
		return Config{}, errors.Wrap(err, "parse postgres env")
	}

	if err := envconfig.Process("", &cnf.http); err != nil {
		return Config{}, errors.Wrap(err, "parse http env")
	}

	if err := envconfig.Process("", &cnf.fastForexAPI); err != nil {
		return Config{}, errors.Wrap(err, "parse fast forex api env")
	}

	return cnf, nil
}

func (c Config) AppName() string {
	return c.app.Name
}

func (c Config) LogLevel() zerolog.Level {
	return zerolog.Level(c.app.LogLevel)
}

func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.http.Port)
}
