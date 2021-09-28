package dislog

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/log"
	"github.com/sirupsen/logrus"
)

var DefaultConfig = Config{
	Logger:    log.Default(),
	LogLevels: ErrorLevelAndAbove,
}

type Config struct {
	Logger        log.Logger
	LogLevels     []logrus.Level
	WebhookID     discord.Snowflake
	WebhookToken  string
	WebhookClient *webhook.Client
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
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

func WithWebhookIDToken(webhookID discord.Snowflake, webhookToken string) ConfigOpt {
	return func(config *Config) {
		config.WebhookID = webhookID
		config.WebhookToken = webhookToken
	}
}

func WithWebhookClient(webhookClient *webhook.Client) ConfigOpt {
	return func(config *Config) {
		config.WebhookClient = webhookClient
	}
}
