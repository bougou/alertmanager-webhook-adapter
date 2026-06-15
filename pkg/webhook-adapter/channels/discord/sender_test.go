package discord

import (
	"fmt"
	"os"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/bwmarrin/discordgo"
)

func Test_Webhook_Send(t *testing.T) {
	id := os.Getenv("DISCORD_WEBHOOK_ID")
	token := os.Getenv("DISCORD_WEBHOOK_TOKEN")

	ws := NewWebhookSender(id, token)

	payload := &models.Payload{
		Title:    "hello",
		Markdown: "this `supports` __a__ **subset** *of* ~~markdown~~ 😃 ```js\nfunction foo(bar) {\n  console.log(bar);\n}\n\nfoo(1);```",
		Text:     "",
	}
	if err := ws.Send(payload); err != nil {
		t.Error(err)
	}
}

func Test_Webhook_Send2(t *testing.T) {
	id := os.Getenv("DISCORD_WEBHOOK_ID")
	token := os.Getenv("DISCORD_WEBHOOK_TOKEN")

	fmt.Println(id, token)
	ws := NewWebhookSender(id, token)

	msg := &MsgWebhook{
		discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Footer: &discordgo.MessageEmbedFooter{
						Text: "this is footer text",
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Alert Name",
							Value:  "yy",
							Inline: true,
						},
						{
							Name:   "Alert Severity",
							Value:  "yy",
							Inline: true,
						},
						{
							Name:  "\t",
							Value: "\t",
						},
						{
							Name:   "Alert Instance",
							Value:  "yy",
							Inline: true,
						},
						{
							Name:   "Alert Instance2",
							Value:  "yy",
							Inline: true,
						},
					},
				},
			},
		},
	}
	if err := ws.SendMsg(msg); err != nil {
		t.Error(err)
	}
}
