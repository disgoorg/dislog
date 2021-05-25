package dislog

import (
	"errors"
	"net/http"
	"strings"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/DisgoOrg/log"
	"github.com/sirupsen/logrus"
)

var ErrNoWebhookSet = errors.New("no webhookClient id token set")

func NewDisLogBuilder() *Builder {
	return &Builder{}
}

type Builder struct {
	httpClient      *http.Client
	logger          log.Logger
	webhookClient   api.WebhookClient
	webhookLogLevel *logrus.Level
	webhookID       *api.Snowflake
	webhookToken    *string
	webhookIDToken  *string
	levels          []logrus.Level
}

func (b *Builder) SetHttpClient(httpClient *http.Client) *Builder {
	b.httpClient = httpClient
	return b
}

func (b *Builder) SetLogger(logger log.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) SetWebhookClient(webhookClient api.WebhookClient) *Builder {
	b.webhookClient = webhookClient
	return b
}

func (b *Builder) SetWebhookLoglevel(webhookLogLevel logrus.Level) *Builder {
	b.webhookLogLevel = &webhookLogLevel
	return b
}

func (b *Builder) SetWebhookIDToken(webhookIDToken string) *Builder {
	b.webhookIDToken = &webhookIDToken
	return b
}

func (b *Builder) SetSetWebhookID(webhookID string) *Builder {
	snowflake := api.Snowflake(webhookID)
	b.webhookID = &snowflake
	return b
}

func (b *Builder) SetSetWebhookIDSnowflake(webhookID api.Snowflake) *Builder {
	b.webhookID = &webhookID
	return b
}

func (b *Builder) SetSetWebhookToken(webhookToken string) *Builder {
	b.webhookToken = &webhookToken
	return b
}

func (b *Builder) SetLevels(levels ...logrus.Level) *Builder {
	b.levels = levels
	return b
}

func (b *Builder) Build() (*DisLog, error) {
	dlog := &DisLog{}

	if b.levels == nil {
		b.levels = WarnLevelAndAbove
	}
	dlog.levels = b.levels

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	if b.webhookIDToken != nil {
		webhookTokenSplit := strings.SplitN(*b.webhookIDToken, "/", 2)
		if len(webhookTokenSplit) != 2 {
			return nil, api.ErrMalformedWebhookToken
		}
		snowflake := api.Snowflake(webhookTokenSplit[0])
		b.webhookID = &snowflake
		b.webhookToken = &webhookTokenSplit[1]
	}

	if b.logger == nil {
		b.logger = logrus.New()
	}

	if b.webhookLogLevel == nil {
		level := logrus.ErrorLevel
		b.webhookLogLevel = &level
	}

	if b.webhookClient == nil {
		if b.webhookID == nil || b.webhookToken == nil {
			return nil, ErrNoWebhookSet
		}
		var err error
		b.webhookClient, err = disgohook.NewWebhookClientByIDToken(b.httpClient, b.logger, *b.webhookID, *b.webhookToken)
		if err != nil {
			return nil, err
		}
	}
	dlog.webhookClient = b.webhookClient

	return dlog, nil
}
