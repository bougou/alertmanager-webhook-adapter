package slack

import (
	"fmt"
	"os"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func Test_SlackSender(t *testing.T) {
	token := os.Getenv("SLACK_APP_TOKEN")
	channel := "#jenkins-ci"
	fmt.Println("slack token:", token)
	sender := NewSender(token, channel, MsgTypeMarkdown)

	payload := &models.Payload{
		Title:    "Hello",
		Text:     "test",
		Markdown: "*Hello*",
	}
	sender.Send(payload)
}
