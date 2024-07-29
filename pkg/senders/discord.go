package senders

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/discord"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeDiscordWebhook = "discord-webhook"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeDiscordWebhook, createDiscordWebhookSender)
}

func createDiscordWebhookSender(request *restful.Request) (models.Sender, error) {
	id := request.QueryParameter("id")
	if id == "" {
		return nil, fmt.Errorf("not id found for discord-webhook channel")
	}

	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for discord-webhook channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender models.Sender = discord.NewWebhookSender(id, token)
	return sender, nil
}
