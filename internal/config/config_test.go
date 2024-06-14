package config_test

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veleton777/test_work_blum/internal/config"
	"os"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	env := map[string]string{
		"APP_NAME":  "app",
		"LOG_LEVEL": "3",

		"HTTP_PORT": "9090",

		"PG_HOST":     "postgres",
		"PG_DATABASE": "pg-db",
		"PG_USER":     "pg-user",
		"PG_PASSWORD": "pg-password",
		"PG_PORT":     "5432",
		"PG_TIMEOUT":  "5s",

		"FAST_FOREX_API_HOST":     "fast-forex",
		"FAST_FOREX_API_KEY":      "fast-forex-api-key",
		"FAST_FOREX_TASK_DELAY":   "2m",
		"FAST_FOREX_HTTP_TIMEOUT": "3s",
	}

	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			require.NoError(t, err)
		}
	}

	conf, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, conf.AppName(), "app")
	assert.Equal(t, conf.LogLevel(), zerolog.Level(3))
	assert.Equal(t, conf.HTTPAddr(), ":9090")
	assert.Equal(t, conf.PgHost(), "postgres")
	assert.Equal(t, conf.PgDB(), "pg-db")
	assert.Equal(t, conf.PgUser(), "pg-user")
	assert.Equal(t, conf.PgPassword(), "pg-password")
	assert.Equal(t, conf.PgPort(), 5432)
	assert.Equal(t, conf.PgTimeout(), time.Second*5)
	assert.Equal(t, conf.FastForexAPIHost(), "fast-forex")
	assert.Equal(t, conf.FastForexAPIKey(), "fast-forex-api-key")
	assert.Equal(t, conf.FastForexBackgroundTaskDelay(), 2*time.Minute)
	assert.Equal(t, conf.FastForexHTTPTimeout(), 3*time.Second)
}
