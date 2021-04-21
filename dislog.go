package dislog

import (
	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = (*DisLog)(nil)

func NewDisLogBuilder() *DisLogBuilder {
	return &DisLogBuilder{}
}

func NewDisLogByToken(webhookToken string, levels ...logrus.Level) (*DisLog, error) {
	dgoHook, err := disgohook.NewDisgoHookByToken(webhookToken)
	if err != nil {
		return nil, err
	}
	return &DisLog{
		disgoHook: dgoHook,
		levels:    levels,
	}, nil
}

func NewDisLogByID(webhookID string, token string, levels ...logrus.Level) (*DisLog, error) {
	dgoHook, err := disgohook.NewDisgoHookByID(webhookID, token)
	if err != nil {
		return nil, err
	}
	return &DisLog{
		disgoHook: dgoHook,
		levels:    levels,
	}, nil
}

type DisLog struct {
	disgoHook api.DisgoHook
	levels    []logrus.Level
}

func (l *DisLog) Levels() []logrus.Level {
	return l.levels
}

func (l *DisLog) Fire(entry *logrus.Entry) error {
	_, err := l.disgoHook.SendMessage()
	if err != nil {
		return err
	}
	return nil
}
