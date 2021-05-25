package main

import (
	"net/http"
	"os"
	"time"

	"github.com/DisgoOrg/dislog"
	"github.com/sirupsen/logrus"
)

var webhookToken = os.Getenv("webhook_token")

var logger = logrus.New()
var httpClient = http.DefaultClient

func main() {
	// override default trace color
	dislog.TraceLevelColor = 0xd400ff

	// override default max delay
	dislog.LogWait = 1 * time.Second

	// override default time format
	dislog.TimeFormatter = "2006-01-02 15:04:05 Z07"

	logger.SetLevel(logrus.TraceLevel)
	logger.Info("starting examples...")
	dlog, err := dislog.NewDisLogBuilder().
		// Sets which logrus.LogLevel the webhook client should log at if none is passed
		SetWebhookLoglevel(logrus.ErrorLevel).
		SetWebhookIDToken(webhookToken).
		// Sets a custom http client or nil for inbuilt
		SetHttpClient(httpClient).
		// Sets which logging levels to send to the webhook
		SetLevels(dislog.TraceLevelAndAbove...).
		Build()
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
	// Calls os.Exit(1) after logging
	logger.Fatal("fatal log")
	// Calls panic() after logging
	logger.Panic("panic log")
}
