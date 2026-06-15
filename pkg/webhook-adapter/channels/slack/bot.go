package slack

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack"
)

// SlackBot can send messages to slack channel
type SlackBot struct {
	token   string
	channel string

	client *slack.Client
}

func NewSlackBot(token string, channel string) *SlackBot {
	client := slack.New(token, slack.OptionDebug(false))
	return &SlackBot{
		token:   token,
		channel: channel,
		client:  client,
	}
}

func (s *SlackBot) Send(msg Msg) error {
	if _, err := s.client.AuthTest(); err != nil {
		msg := fmt.Sprintf("slack auth failed, err: %s", err)
		return errors.New(msg)
	}

	_, _, err := s.client.PostMessage(
		s.channel,
		slack.MsgOptionBlocks(msg...),
	)
	if err != nil {
		msg := fmt.Sprintf("slack send failed, err: %s", err)
		return errors.New(msg)
	}
	return nil
}
