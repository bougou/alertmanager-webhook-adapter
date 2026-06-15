package slack

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/slack-go/slack"
)

func init() {
	Payload2MsgFnMap[MsgTypeText] = NewMsgTextFromPayload

}
func NewMsgTextFromPayload(payload *models.Payload) Msg {
	headerText := slack.NewTextBlockObject(
		slack.PlainTextType,
		payload.Title,
		true, false,
	)
	if err := headerText.Validate(); err != nil {
		return nil
	}
	headerBlock := slack.NewHeaderBlock(headerText)

	bodyText := slack.NewTextBlockObject(slack.PlainTextType, payload.Text, false, false)
	if err := bodyText.Validate(); err != nil {
		return nil
	}
	bodyBlock := slack.NewSectionBlock(bodyText, nil, nil)

	return []slack.Block{headerBlock, bodyBlock}
}
