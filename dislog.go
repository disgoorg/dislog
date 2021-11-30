package dislog

import (
	"errors"
	"fmt"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/log"
	"strconv"
	"sync"
	"time"

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
	config := &DefaultConfig
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}

	if config.WebhookClient == nil {
		if config.WebhookID == "" || config.WebhookToken == "" {
			return nil, errors.New("please provide either a webhook client or id & token")
		}
		config.WebhookClient = webhook.NewClient(config.WebhookID, config.WebhookToken, webhook.WithLogger(config.Logger))
	}
	return &DisLog{
		webhookClient: config.WebhookClient,
		levels:        config.LogLevels,
	}, nil
}

type DisLog struct {
	sync.Mutex

	webhookClient *webhook.Client
	queued        bool
	queue         []discord.Embed
	levels        []logrus.Level
}

func (l *DisLog) Levels() []logrus.Level {
	return l.levels
}

func (l *DisLog) Close() {
	l.sendEmbeds()
}

func (l *DisLog) queueEmbed(embed discord.Embed, forceSend bool) error {
	l.Lock()
	defer l.Unlock()
	l.queue = append(l.queue, embed)
	if len(l.queue) >= MaxEmbeds || forceSend {
		go l.sendEmbeds()
	} else {
		l.queueEmbeds()
	}
	return nil
}

func (l *DisLog) sendEmbeds() {
	l.Lock()
	defer l.Unlock()
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

	_, err := l.webhookClient.CreateMessage(message.Build())
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
