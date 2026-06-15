package models

type Sender interface {
	Send(payload *Payload) error
	SendMsg(msg interface{}) error
	SendMsgT(msgType string, msg interface{}) error
}
