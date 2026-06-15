package dingtalk

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeText] = NewMsgTextFromPayload
}

type Text struct {
	Content string `json:"content"`
}

func NewText(content string) *Text {
	return &Text{content}
}

func NewMsgText(text *Text) *Msg {
	return &Msg{
		MsgType: MsgTypeText,
		Text:    text,
	}
}

func (bot *DingtalkGroupBot) SendText(content string, atMobiles []string, atAll bool) error {
	text := NewText(content)
	msg := NewMsgText(text)
	msg.WithAtAll(atAll).WithAtMobiles(atMobiles)
	return bot.Send(msg)
}

func NewMsgTextFromPayload(payload *models.Payload) *Msg {
	return &Msg{
		MsgType: MsgTypeText,
		Text:    NewText(payload.Text),
	}
}
