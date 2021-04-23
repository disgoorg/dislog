package dislog

import (
	"fmt"
	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	PanicLevelColor = 0xe300bd
	FatalLevelColor = 0xeb4034
	ErrorLevelColor = 0xff0011
	WarnLevelColor  = 0xfcec00
	InfoLevelColor  = 0x002ccc
	DebugLevelColor = 0x0095ff
	TraceLevelColor = 0xffffff
)

var LevelColors = map[logrus.Level]int{
	logrus.PanicLevel: PanicLevelColor,
	logrus.FatalLevel: FatalLevelColor,
	logrus.ErrorLevel: ErrorLevelColor,
	logrus.WarnLevel:  WarnLevelColor,
	logrus.InfoLevel:  InfoLevelColor,
	logrus.DebugLevel: DebugLevelColor,
	logrus.TraceLevel: TraceLevelColor,
}

var _ logrus.Hook = (*DisLog)(nil)

func NewDisLogBuilder() *DisLogBuilder {
	return &DisLogBuilder{}
}

func NewDisLogByToken(webhookToken string, levels ...logrus.Level) (*DisLog, error) {
	webhook, err := disgohook.NewWebhookByToken(webhookToken)
	if err != nil {
		return nil, err
	}
	return &DisLog{
		webhook: webhook,
		levels:  levels,
	}, nil
}

func NewDisLogByID(webhookID string, token string, levels ...logrus.Level) (*DisLog, error) {
	webhook, err := disgohook.NewWebhookByID(webhookID, token)
	if err != nil {
		return nil, err
	}
	return &DisLog{
		webhook: webhook,
		levels:  levels,
	}, nil
}

type DisLog struct {
	webhook api.Webhook
	levels  []logrus.Level
}

func (l *DisLog) Levels() []logrus.Level {
	return l.levels
}

func (l *DisLog) Fire(entry *logrus.Entry) error {
	embed := api.NewEmbedBuilder().
		SetColor(LevelColors[entry.Level]).
		SetTitle("log in file: "+entry.Caller.File+" from func: "+entry.Caller.Function+" in line: "+strconv.Itoa(entry.Caller.Line)).
		SetDescriptionf(entry.Message).
		AddField("Level", entry.Level.String(), true).
		AddField("Time", entry.Time.String(), true)
	for key, value := range entry.Data {
		embed.AddField(key, fmt.Sprint(value), true)
	}
	_, err := l.webhook.SendMessage(api.NewWebhookMessageWithEmbeds(embed.Build()).Build())
	if err != nil {
		return err
	}
	return nil
}

func subStr(str string, start int, end int) string {
	a := []rune(str)
	return string(a[start:end])
}
