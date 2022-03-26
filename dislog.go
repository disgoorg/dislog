package dislog

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
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

func New(opts ...ConfigOpt) (*DisLog, error) {
	config := DefaultConfig()
	config.Apply(opts)

	if config.WebhookID == "" || config.WebhookToken == "" {
		return nil, errors.New("Webhook ID & token or a webhook client are required")
	}

	return &DisLog{
		config: *config,
	}, nil
}

type DisLog struct {
	config Config

	queued  bool
	queue   []discord.Embed
	queueMu sync.Mutex
}

func (l *DisLog) Levels() []logrus.Level {
	return l.config.LogLevels
}

func (l *DisLog) Close(ctx context.Context) {
	l.sendEmbeds()
	l.config.WebhookClient.Close(ctx)
}

func (l *DisLog) queueEmbed(embed discord.Embed, forceSend bool) error {
	l.queueMu.Lock()
	defer l.queueMu.Unlock()
	l.queue = append(l.queue, embed)
	if len(l.queue) >= MaxEmbeds || forceSend {
		go l.sendEmbeds()
	} else {
		l.queueEmbeds()
	}
	return nil
}

func (l *DisLog) sendEmbeds() {
	l.queueMu.Lock()
	defer l.queueMu.Unlock()
	if len(l.queue) == 0 {
		return
	}
	message := discord.NewWebhookMessageCreateBuilder()

	for i := 0; i < len(l.queue); i++ {
		if i >= MaxEmbeds {
			// queue again as we have logs to send
			l.queueEmbeds()
			break
		}
		message.AddEmbeds(l.queue[i])
		l.queue = append(l.queue[:i], l.queue[i+1:]...)
		i--
	}
	if len(message.Embeds) == 0 {
		return
	}

	_, err := l.config.WebhookClient.CreateMessage(message.Build())
	if err != nil {
		fmt.Printf("error while sending logs: %s\n", err)
	}
}

func (l *DisLog) queueEmbeds() {
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
	eb := discord.NewEmbedBuilder().
		SetColor(*LevelColors[entry.Level]).
		SetDescription(entry.Message).
		AddField("LogLevel", entry.Level.String(), true).
		AddField("Time", entry.Time.Format(TimeFormatter), true)
	if entry.HasCaller() {
		eb.SetDescription("in file: " + entry.Caller.File + " from func: " + entry.Caller.Function + " in line: " + strconv.Itoa(entry.Caller.Line) + "\n" + entry.Message)
	}
	for key, value := range entry.Data {
		eb.AddField(key, fmt.Sprint(value), true)
	}
	return l.queueEmbed(eb.Build(), entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel)
}
