package dingtalk

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeMarkdown] = NewMsgMarkdownFromPayload
}

type Markdown struct {
	Title string `json:"title"` // 首屏会话透出的展示内容
	Text  string `json:"text"`
}

func (md *Markdown) Valid() bool {
	if md.Title == "" || md.Text == "" {
		return false
	}
	return true
}

func NewMarkdown(title string, text string) *Markdown {
	// DingTalk Markdown Title does not show on the Text page
	t := SanitizeMarkdown(fmt.Sprintf("%s\n%s", title, text))

	return &Markdown{
		Title: title,
		Text:  t,
	}
}

func NewMsgMarkdown(md *Markdown) *Msg {
	return &Msg{
		MsgType:  MsgTypeMarkdown,
		Markdown: md,
	}
}

func NewMsgMarkdownFromPayload(payload *models.Payload) *Msg {
	md := NewMarkdown(payload.Title, payload.Markdown)
	msg := NewMsgMarkdown(md)
	msg.WithAtMobiles(payload.At.AtMobiles)
	msg.WithAtAll(payload.At.AtAll)

	return msg
}

func SanitizeMarkdown(content string) string {
	return content
}
