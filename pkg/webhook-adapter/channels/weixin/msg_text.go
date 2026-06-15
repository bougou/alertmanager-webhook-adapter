package weixin

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/utils"
)

func init() {
	Payload2MsgFnMap[MsgTypeText] = NewMsgTextFromPayload
}

type Text struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func (text *Text) WithMentionedList(mentionedList []string) *Text {
	text.MentionedList = mentionedList
	return text
}

func (text *Text) WithMentionedMobileList(mentionedMobileList []string) *Text {
	text.MentionedMobileList = mentionedMobileList
	return text
}

type TextOption func(*Text)

func TextMentionedListOpt(mentionedList []string) TextOption {
	return func(text *Text) {
		text.MentionedList = mentionedList
	}
}

func TextMentionedMobileListOpt(mentionedMobileList []string) TextOption {
	return func(text *Text) {
		text.MentionedMobileList = mentionedMobileList
	}
}

func NewText(content string, options ...TextOption) *Text {
	text := &Text{
		Content: utils.TruncateToValidUTF8(content, maxTextBytes, truncatedMark),
	}

	for _, option := range options {
		option(text)
	}

	return text
}

func NewMsgText(text *Text) *Msg {
	return &Msg{
		MsgType: MsgTypeText,
		Text:    text,
	}
}

func NewMsgTextFromPayload(payload *models.Payload) *Msg {
	text := NewText(payload.Text)
	text.WithMentionedMobileList(payload.At.AtMobiles)

	return NewMsgText(text)
}
