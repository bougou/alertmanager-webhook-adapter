package weixinapp

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeVoice] = NewMsgVoiceFromPayload
}

type Voice struct {
	MediaID string `json:"media_id"`
}

// Todo
func NewMsgVoiceFromPayload(payload *models.Payload) *Msg {
	return &Msg{}
}
