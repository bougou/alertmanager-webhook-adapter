package senders

import (
	"fmt"
	"strings"

	"github.com/bougou/webhook-adapter/channels/slack"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeSlack = "slack"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeSlack, createSlackSender)
}

func createSlackSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for slack")
	}
	channel := request.QueryParameter("channel")
	if channel == "" {
		return nil, fmt.Errorf("not channel found for slack")
	}
	// add # if channel not begin with #
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender models.Sender = slack.NewSender(token, channel, msgType)
	return sender, nil
}
