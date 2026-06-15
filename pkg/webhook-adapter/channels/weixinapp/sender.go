package weixinapp

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

type Sender struct {
	notifer *Notifier
	msgType string
}

var _ models.Sender = (*Sender)(nil)

func NewSender(corpID string, agentID int, agentSecret string, msgType string, toUser string, toParty string, toTag string) models.Sender {
	return newSender(corpID, agentID, agentSecret, msgType, toUser, toParty, toTag)
}

func newSender(corpID string, agentID int, agentSecret string, msgType string, toUser string, toParty string, toTag string) *Sender {
	return &Sender{
		notifer: NewNotifer(corpID, agentID, agentSecret, toUser, toParty, toTag),
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
		return fmt.Errorf("passed msgSource is not type *weixinapp.Msg")
	}

	return s.notifer.Send(msg)
}
