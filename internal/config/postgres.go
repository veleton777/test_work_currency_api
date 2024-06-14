package config

import "time"

type postgres struct {
	DB       string        `envconfig:"PG_DATABASE"`
	Host     string        `envconfig:"PG_HOST"`
	Port     int           `envconfig:"PG_PORT"`
	User     string        `envconfig:"PG_USER"`
	Password string        `envconfig:"PG_PASSWORD"`
	Timeout  time.Duration `envconfig:"PG_TIMEOUT"`
}

func (c Config) PgHost() string {
	return c.postgres.Host
}

func (c Config) PgDB() string {
	return c.postgres.DB
}

func (c Config) PgUser() string {
	return c.postgres.User
}

func (c Config) PgPassword() string {
	return c.postgres.Password
}

func (c Config) PgPort() int {
	return c.postgres.Port
}

func (c Config) PgTimeout() time.Duration {
	return c.postgres.Timeout
}
