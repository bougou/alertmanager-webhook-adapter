package weixinapp

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeText] = NewMsgTextFromPayload
}

type Text struct {
	Content string `json:"content"`
}

func NewMsgText(content string) *Msg {
	return &Msg{
		MsgType: MsgTypeText,
		Text: &Text{
			Content: content,
		},
	}
}

func NewMsgTextFromPayload(payload *models.Payload) *Msg {
	return NewMsgText(payload.Text)
}
