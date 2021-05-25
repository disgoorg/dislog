package dislog

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/DisgoOrg/log"
	"github.com/sirupsen/logrus"
)

const MaxEmbeds = 10

var (
	TimeFormatter = "2006-01-02 15:04:05 Z07"
	LogWait       = 2 * time.Second

	PanicLevelColor = 0xe300bd
	FatalLevelColor = 0xff0011
	ErrorLevelColor = 0xeb4034
	WarnLevelColor  = 0xff8c00
	InfoLevelColor  = 0xfcec00
	DebugLevelColor = 0x0095ff
	TraceLevelColor = 0xfffffe

	LevelColors = map[logrus.Level]*int{
		logrus.PanicLevel: &PanicLevelColor,
		logrus.FatalLevel: &FatalLevelColor,
		logrus.ErrorLevel: &ErrorLevelColor,
		logrus.WarnLevel:  &WarnLevelColor,
		logrus.InfoLevel:  &InfoLevelColor,
		logrus.DebugLevel: &DebugLevelColor,
		logrus.TraceLevel: &TraceLevelColor,
	}

	PanicLevelAndAbove = []logrus.Level{logrus.PanicLevel}
	FatalLevelAndAbove = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	ErrorLevelAndAbove = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	WarnLevelAndAbove  = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
	InfoLevelAndAbove  = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}
	DebugLevelAndAbove = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel}
	TraceLevelAndAbove = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}

	_ logrus.Hook = (*DisLog)(nil)
)

func newDisLog(webhookLogLevel logrus.Level, webhookClientCreate func(logger log.Logger) (api.WebhookClient, error), levels ...logrus.Level) (*DisLog, error) {
	logger := logrus.New()
	logger.SetLevel(webhookLogLevel)
	webhook, err := webhookClientCreate(logger)
	if err != nil {
		return nil, err
	}
	return NewDisLog(webhook, levels...)
}

func NewDisLogByToken(httpClient *http.Client, webhookLogLevel logrus.Level, webhookToken string, levels ...logrus.Level) (*DisLog, error) {
	return newDisLog(webhookLogLevel, func(logger log.Logger) (api.WebhookClient, error) {
		return disgohook.NewWebhookClientByToken(httpClient, logger, webhookToken)
	}, levels...)
}

func NewDisLogByIDSnowflakeToken(httpClient *http.Client, webhookLogLevel logrus.Level, webhookID api.Snowflake, webhookToken string, levels ...logrus.Level) (*DisLog, error) {
	return newDisLog(webhookLogLevel, func(logger log.Logger) (api.WebhookClient, error) {
		return disgohook.NewWebhookClientByIDToken(httpClient, logger, webhookID, webhookToken)
	}, levels...)
}

func NewDisLogByIDToken(httpClient *http.Client, webhookLogLevel logrus.Level, webhookID string, webhookToken string, levels ...logrus.Level) (*DisLog, error) {
	return NewDisLogByIDSnowflakeToken(httpClient, webhookLogLevel, api.Snowflake(webhookID), webhookToken, levels...)
}

func NewDisLog(webhook api.WebhookClient, levels ...logrus.Level) (*DisLog, error) {
	return &DisLog{
		webhookClient: webhook,
		queued:        false,
		levels:        levels,
	}, nil
}

type DisLog struct {
	webhookClient api.WebhookClient
	lock          sync.Mutex
	queued        bool
	logQueue      []*api.Embed
	levels        []logrus.Level
}

func (l *DisLog) Levels() []logrus.Level {
	return l.levels
}

func (l *DisLog) Close() {
	l.sendEmbeds()
}

func (l *DisLog) queueEmbed(embed *api.Embed, forceSend bool) error {
	l.logQueue = append(l.logQueue, embed)
	if len(l.logQueue) >= MaxEmbeds || forceSend {
		l.sendEmbeds()
	} else {
		l.queueSendEmbeds()
	}
	return nil
}

func (l *DisLog) sendEmbeds() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.logQueue) == 0 {
		return
	}
	message := api.NewWebhookMessageBuilder()

	for i := 0; i < len(l.logQueue); i++ {
		if i >= MaxEmbeds {
			// queue again as we have logs to send
			l.queueSendEmbeds()
			break
		}
		message.AddEmbeds(l.logQueue[i])
		l.logQueue = append(l.logQueue[:i], l.logQueue[i+1:]...)
		i--
	}
	if len(message.Embeds) == 0 {
		return
	}

	_, err := l.webhookClient.SendMessage(message.Build())
	if err != nil {
		fmt.Printf("error while sending logs: %s\n", err)
	}
}

func (l *DisLog) queueSendEmbeds() {
	if l.queued {
		return
	}
	go func() {
		l.queued = true
		time.Sleep(LogWait)
		l.sendEmbeds()
		l.queued = false
	}()
}

func (l *DisLog) Fire(entry *logrus.Entry) error {
	eb := api.NewEmbedBuilder().
		SetColor(*LevelColors[entry.Level]).
		SetDescription(entry.Message).
		AddField("Level", entry.Level.String(), true).
		AddField("Time", entry.Time.Format(TimeFormatter), true)
	if entry.HasCaller() {
		eb.SetDescription("in file: " + entry.Caller.File + " from func: " + entry.Caller.Function + " in line: " + strconv.Itoa(entry.Caller.Line) + "\n" + entry.Message)
	}
	for key, value := range entry.Data {
		eb.AddField(key, fmt.Sprint(value), true)
	}
	return l.queueEmbed(eb.Build(), entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel)
}
