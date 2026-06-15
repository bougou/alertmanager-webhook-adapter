package feishu

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

const (
	ChannelTypeFeishu string = "feishu"

	MsgTypeImage       string = "image"
	MsgTypeInteractive string = "interactive"
	MsgTypeMarkdown    string = "markdown" // Underlying, we use interactive msg type to implement markdown
	MsgTypePost        string = "post"
	MsgTypeShareChat   string = "sharechat"
	MsgTypeText        string = "text"
)

type Msg struct {
	// 开启签名验证后发送文本消息
	// Timestamp time.Time `json:"timestamp,omitempty"`
	// Sign      string    `json:"sign,omitempty"`

	MsgType string `json:"msg_type"`

	Content *Content `json:"content,omitempty"`

	Card        *Card  `json:"card,omitempty"`
	RootID      string `json:"root_id,omitempty"`      // 需要回复的消息的open_message_id
	UpdateMulti bool   `json:"update_multi,omitempty"` // 控制卡片是否是共享卡片(所有用户共享同一张消息卡片），默认为 false
}

type Content struct {
	Text        string `json:"text,omitempty"`
	ImageKey    string `json:"image,omitempty"`
	Post        *Post  `json:"post,omitempty"`
	ShareChatID string `json:"share_chat_id,omitempty"`
}

type Payload2MsgFn func(payload *models.Payload) *Msg

var Payload2MsgFnMap = make(map[string]Payload2MsgFn)

func validateMsg(msgType string, msg *Msg) error {
	if msg.MsgType != msgType {

		if msgType == MsgTypeMarkdown {
			// markdown msg currently implemented as interactive msg.
			if msg.MsgType == MsgTypeInteractive || msg.MsgType == MsgTypeMarkdown {
				return nil
			}
		}

		return fmt.Errorf("the msg does not match with specified msgType")
	}

	return nil
}
