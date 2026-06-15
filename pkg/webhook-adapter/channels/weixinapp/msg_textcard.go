package weixinapp

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeTextCard] = NewMsgTextCardFromPayload
}

type TextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BtnText     string `json:"btntext"`
}

// Todo
func NewMsgTextCardFromPayload(payload *models.Payload) *Msg {
	return &Msg{}
}
