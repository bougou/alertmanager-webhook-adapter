package feishu

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/feishu/card"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeMarkdown] = NewMsgMarkdownFromPayload
}

func NewMsgMarkdown(title string, markdown string) *Msg {
	card := NewCardMarkdown(title, markdown, "")
	return &Msg{
		MsgType: MsgTypeInteractive,
		Card:    card,
	}
}

func NewMsgMarkdownFromPayload(payload *models.Payload) *Msg {
	card := NewCardMarkdown(payload.Title, payload.Markdown, headerColor(payload.EffectiveSeverity()))
	return &Msg{
		MsgType: MsgTypeInteractive,
		Card:    card,
	}
}

func NewCardMarkdown(title string, markdown string, headerColor string) *Card {
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

	if headerColor == "" {
		headerColor = "turquoise"
	}

	return &Card{
		Config: &CardConfig{
			EnableForward: false,
		},
		Header: &CardHeader{
			Title: &card.Text{
				Tag:     "plain_text",
				Content: title,
			},
			Template: headerColor,
		},
		Elements: elements,
	}
}

// headerColor maps shared effective severity to Feishu card header colors.
func headerColor(severity string) string {
	switch severity {
	case models.SeverityCritical:
		return "red"
	case models.SeverityWarning:
		return "orange"
	case models.SeverityInfo:
		return "blue"
	case models.SeverityOK:
		return "green"
	default:
		return "turquoise"
	}
}
