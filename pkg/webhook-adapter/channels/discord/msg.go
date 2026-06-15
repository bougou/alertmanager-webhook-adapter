package discord

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/bwmarrin/discordgo"
)

const (
	ChannelTypeDiscordWebhook string = "discord-webhook"
)

type MsgWebhook struct {
	discordgo.WebhookParams `json:"inline"`
}

func NewMsgWebhookFromPayload(payload *models.Payload) *MsgWebhook {
	msg := &MsgWebhook{
		discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       payload.Title,
					Description: payload.Markdown,
				},
			},
		},
	}

	return msg
}
