package senders

import (
	"fmt"
	"strings"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/channels/slack"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeSlack = "slack"
)

type (
	slackMarkdownConverter struct{}
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeSlack, createSlackSender)

	RegisterChannelsMsgConverter(ChannelTypeSlack, slack.MsgTypeMarkdown, &slackMarkdownConverter{})
}

func createSlackSender(request *restful.Request) (sender models.Sender, converter MsgConverter, err error) {
	token := request.QueryParameter("token")
	if token == "" {
		err = fmt.Errorf("not token found for slack")
		return
	}
	channel := request.QueryParameter("channel")
	if channel == "" {
		err = fmt.Errorf("not channel found for slack")
		return
	}
	// add # if channel not begin with #
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	sender = slack.NewSender(token, channel, msgType)

	converter, ok := getMsgConverter(ChannelTypeSlack, msgType)
	if !ok {
		err = ErrNotFoundConverter(ChannelTypeSlack, msgType)

		return
	}

	return
}

func (c *slackMarkdownConverter) Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error) {
	payload, err := promMsg.ToPayload(ChannelTypeSlack, raw)
	if err != nil {
		return nil, fmt.Errorf("create msg payload failed, %v", err)
	}

	msgType := slack.MsgTypeMarkdown
	payload2MsgFn, exists := slack.Payload2MsgFnMap[msgType]
	if !exists {
		return nil, ErrNotFoundPayload2MsgFn(ChannelTypeSlack, msgType)
	}
	msg := payload2MsgFn(payload)

	return msg, nil
}
