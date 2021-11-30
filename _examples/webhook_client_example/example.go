package main

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/webhook"
	"net/http"
	"os"
	"time"

	"github.com/DisgoOrg/dislog"
	"github.com/sirupsen/logrus"
)

var (
	webhookID    = discord.Snowflake(os.Getenv("webhook_id"))
	webhookToken = os.Getenv("webhook_token")
)

func main() {
	logger := logrus.New()
	// override default trace color
	dislog.TraceLevelColor = 0xd400ff

	// override default max delay
	dislog.LogWait = 1 * time.Second

	// override default time format
	dislog.TimeFormatter = "2006-01-02 15:04:05 Z07"

	logger.SetLevel(logrus.TraceLevel)
	logger.Info("starting example...")

	dlog, err := dislog.New(
		// Sets which logging levels to send to the webhook
		dislog.WithLogLevels(dislog.TraceLevelAndAbove...),
		// Sets a custom http client or nil for inbuilt
		dislog.WithWebhookClient(webhook.NewClient(webhookID, webhookToken, webhook.WithRestClientConfigOpts(rest.WithHTTPClient(http.DefaultClient)))),
	)
	if err != nil {
		logger.Errorf("error initializing dislog %s", err)
		return
	}
	defer dlog.Close()
	logger.AddHook(dlog)

	logger.Trace("trace log")
	logger.Debug("debug log")
	logger.Info("info log")
	logger.Warn("warn log")
	logger.Error("error log")
	// Calls panic() after logging
	logger.Panic("panic log")
	// Calls os.Exit(1) after logging
	logger.Fatal("fatal log")
}
