package slack

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/slack-go/slack"
)

const (
	ChannelTypeSlack string = "slack"

	MsgTypeMarkdown string = "markdown"
	MsgTypeText     string = "text"
)

type Payload2MsgFn func(payload *models.Payload) Msg

var Payload2MsgFnMap = make(map[string]Payload2MsgFn)

type Msg []slack.Block

func validateMsg(msgType string, msg *Msg) error {
	return nil
}
