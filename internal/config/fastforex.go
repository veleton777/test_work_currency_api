package config

import "time"

type fastForexAPI struct {
	Host                string        `envconfig:"FAST_FOREX_API_HOST"`
	APIKey              string        `envconfig:"FAST_FOREX_API_KEY"`
	BackgroundTaskDelay time.Duration `envconfig:"FAST_FOREX_TASK_DELAY"`
	HTTPTimeout         time.Duration `envconfig:"FAST_FOREX_HTTP_TIMEOUT"`
}

func (c Config) FastForexAPIHost() string {
	return c.fastForexAPI.Host
}

func (c Config) FastForexAPIKey() string {
	return c.fastForexAPI.APIKey
}

func (c Config) FastForexBackgroundTaskDelay() time.Duration {
	return c.fastForexAPI.BackgroundTaskDelay
}

func (c Config) FastForexHTTPTimeout() time.Duration {
	return c.fastForexAPI.HTTPTimeout
}
