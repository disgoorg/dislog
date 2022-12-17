package main

import (
	"context"
	"os"
	"time"

	"github.com/disgoorg/dislog"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

var (
	webhookID    = snowflake.GetEnv("webhook_id")
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

	hook, err := dislog.New(
		// Sets which logging levels to send to the webhook
		dislog.WithLogLevels(dislog.TraceLevelAndAbove...),
		// Sets webhook id & token
		dislog.WithWebhookIDToken(webhookID, webhookToken),
	)
	if err != nil {
		logger.Errorf("error initializing dislog %s", err)
		return
	}
	defer hook.Close(context.TODO())
	logger.AddHook(hook)

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
