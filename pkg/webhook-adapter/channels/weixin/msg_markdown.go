package weixin

import (
	"fmt"
	"strings"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/utils"
)

func init() {
	Payload2MsgFnMap[MsgTypeMarkdown] = NewMsgMarkdownFromPayload
}

type Markdown struct {
	Content string `json:"content"` // this should be raw markdown string, weixin bot only support a small subset syntax
}

func NewMsgMarkdown(markdown *Markdown) *Msg {
	return &Msg{
		MsgType:  MsgTypeMarkdown,
		Markdown: markdown,
	}
}

func NewMsgMarkdownFromContent(content string) *Msg {
	content = SanitizeMarkdown(content)
	truncated := utils.TruncateToValidUTF8(content, maxMarkdownBytes, truncatedMark)

	msg := &Msg{
		MsgType: MsgTypeMarkdown,
		Markdown: &Markdown{
			Content: truncated,
		},
	}
	return msg
}

func NewMsgMarkdownFromPayload(payload *models.Payload) *Msg {
	m := fmt.Sprintf("%s\n%s", payload.Title, payload.Markdown)
	return NewMsgMarkdownFromContent(m)
}

func SanitizeMarkdown(content string) string {
	// no need <br> for line break
	content = strings.ReplaceAll(content, "<br>", "")

	// remove `>` line
	content = strings.ReplaceAll(content, "\n>\n", "\n")

	return content
}
