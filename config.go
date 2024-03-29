package dislog

import (
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:    log.Default(),
		LogLevels: ErrorLevelAndAbove,
	}
}

type Config struct {
	Logger    log.Logger
	LogLevels []logrus.Level

	WebhookID     snowflake.ID
	WebhookToken  string
	WebhookClient webhook.Client
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.WebhookClient == nil && c.WebhookID != 0 && c.WebhookToken != "" {
		c.WebhookClient = webhook.New(c.WebhookID, c.WebhookToken)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithLogLevels(levels ...logrus.Level) ConfigOpt {
	return func(config *Config) {
		config.LogLevels = levels
	}
}

func WithWebhookIDToken(webhookID snowflake.ID, webhookToken string) ConfigOpt {
	return func(config *Config) {
		config.WebhookID = webhookID
		config.WebhookToken = webhookToken
	}
}

func WithWebhookClient(webhookClient webhook.Client) ConfigOpt {
	return func(config *Config) {
		config.WebhookClient = webhookClient
	}
}
