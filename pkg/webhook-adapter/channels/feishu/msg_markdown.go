package feishu

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/feishu/card"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeMarkdown] = NewMsgMarkdownFromPayload
}

func NewMsgMarkdown(title string, markdown string) *Msg {
	card := NewCardMarkdown(title, markdown)
	return &Msg{
		MsgType: MsgTypeInteractive,
		Card:    card,
	}
}

func NewMsgMarkdownFromPayload(payload *models.Payload) *Msg {
	card := NewCardMarkdown(payload.Title, payload.Markdown)
	return &Msg{
		MsgType: MsgTypeInteractive,
		Card:    card,
	}
}

func NewCardMarkdown(title string, markdown string) *Card {
	elements := []card.CardModule{}

	// see: https://open.feishu.cn/document/ukTMukTMukTM/uADOwUjLwgDM14CM4ATN

	// module := &card.ModuleDiv{
	// 	Tag: "div",
	// 	Text: &card.Text{
	// 		Tag:     "lark_md",
	// 		Content: markdown,
	// 	},
	// }

	module := card.NewModuleMarkdown(markdown, nil)

	elements = append(elements, module)

	return &Card{
		Config: &CardConfig{
			EnableForward: false,
		},
		Header: &CardHeader{
			Title: &card.Text{
				Tag:     "plain_text",
				Content: title,
			},
		},
		Elements: elements,
	}
}
