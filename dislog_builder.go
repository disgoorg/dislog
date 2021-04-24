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

var ErrNoWebhookSet = errors.New("no webhook id token set")

func NewDisLogBuilder() *DisLogBuilder {
	return &DisLogBuilder{}
}

type DisLogBuilder struct {
	httpClient     *http.Client
	logger         log.Logger
	webhook        api.Webhook
	webhookID      *string
	webhookToken   *string
	webhookIDToken *string
	levels         []logrus.Level
}

func (b *DisLogBuilder) SetHttpClient(httpClient *http.Client) *DisLogBuilder {
	b.httpClient = httpClient
	return b
}

func (b *DisLogBuilder) SetLogger(logger log.Logger) *DisLogBuilder {
	b.logger = logger
	return b
}

func (b *DisLogBuilder) SetWebhook(webhook api.Webhook) *DisLogBuilder {
	b.webhook = webhook
	return b
}

func (b *DisLogBuilder) SetWebhookIDToken(webhookIDToken string) *DisLogBuilder {
	b.webhookIDToken = &webhookIDToken
	return b
}

func (b *DisLogBuilder) SetSetWebhookID(webhookID string) *DisLogBuilder {
	b.webhookID = &webhookID
	return b
}

func (b *DisLogBuilder) SetSetWebhookToken(webhookToken string) *DisLogBuilder {
	b.webhookToken = &webhookToken
	return b
}

func (b *DisLogBuilder) SetLevels(levels ...logrus.Level) *DisLogBuilder {
	b.levels = levels
	return b
}

func (b *DisLogBuilder) Build() (*DisLog, error) {
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
		b.webhookID = &webhookTokenSplit[0]
		b.webhookToken = &webhookTokenSplit[1]
	}

	if b.webhook == nil {
		if b.webhookID == nil || b.webhookToken == nil {
			return nil, ErrNoWebhookSet
		}
		webhook, err := disgohook.NewWebhookByIDToken(b.httpClient, b.logger, *b.webhookID, *b.webhookToken)
		if err != nil {
			return nil, err
		}
		b.webhook = webhook
	}

	return dlog, nil
}
