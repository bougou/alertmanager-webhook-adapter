package dingtalk

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

const (
	ChannelTypeDingtalk string = "dingtalk"

	MsgTypeActionCard string = "actioncard"
	MsgTypeFeedCard   string = "feedcard"
	MsgTypeLink       string = "link"
	MsgTypeMarkdown   string = "markdown"
	MsgTypeText       string = "text"
)

type Payload2MsgFn func(payload *models.Payload) *Msg

var Payload2MsgFnMap = make(map[string]Payload2MsgFn)

type Msg struct {
	MsgType    string      `json:"msgtype"`
	Text       *Text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	ActionCard *ActionCard `json:"actionCard,omitempty"`
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`
	At         *At         `json:"at,omitempty"` // only available for text and markdown type
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func (msg *Msg) SupportAt() bool {
	if msg.MsgType == MsgTypeText || msg.MsgType == MsgTypeMarkdown {
		return true
	}
	return false
}

func (msg *Msg) WithAt(at *At) *Msg {
	if msg.SupportAt() {
		msg.At = at
	}
	return msg
}

func (msg *Msg) WithAtAll(atAll bool) *Msg {
	if msg.SupportAt() && msg.At == nil {
		msg.At = &At{}
	}
	msg.At.IsAtAll = atAll
	return msg
}

func (msg *Msg) WithAtMobiles(mobiles []string) *Msg {
	if msg.SupportAt() && msg.At == nil {
		msg.At = &At{}
	}
	msg.At.AtMobiles = mobiles
	return msg
}

func validateMsg(msgType string, msg *Msg) error {
	if msg.MsgType != msgType {
		return fmt.Errorf("the msg does not match with specified msgType")
	}
	return nil
}
