package weixin

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

type Sender struct {
	bot     *WeixinGroupBot
	msgType string
}

var _ models.Sender = (*Sender)(nil)

func NewSender(key string, msgType string) models.Sender {
	return newSender(key, msgType)
}

func newSender(key string, msgType string) *Sender {
	if msgType == "" {
		msgType = MsgTypeMarkdown
	}

	return &Sender{
		bot:     NewWexinGroupBot(key),
		msgType: msgType,
	}
}

func (s *Sender) Send(payload *models.Payload) error {
	payload2Msg, ok := Payload2MsgFnMap[s.msgType]
	if !ok {
		return fmt.Errorf("unkown msg type")
	}
	msg := payload2Msg(payload)

	return s.SendMsg(msg)
}

func (s *Sender) SendMsg(msgSource interface{}) error {
	return s.SendMsgT(s.msgType, msgSource)
}

func (s *Sender) SendMsgT(msgType string, msgSource interface{}) error {
	msg, ok := msgSource.(*Msg)
	if !ok {
		return fmt.Errorf("passed msgSource is not type *weixin.Msg")
	}

	switch msgType {
	case MsgTypeFile:
	case MsgTypeImage:
	case MsgTypeMarkdown:
	case MsgTypeNews:
	case MsgTypeText:
	case MsgTypeTemplateCard:
	default:
		return fmt.Errorf("unsupported msgtype of (%s)", msgType)
	}

	if err := validateMsg(msgType, msg); err != nil {
		return fmt.Errorf("valid msg failed, err: %s", err)
	}

	return s.bot.Send(msg)
}
