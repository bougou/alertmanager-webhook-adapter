package weixinapp

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeInteractiveCard] = NewMsgInteractiveCardFromPayload
}

type InteractiveTaskcard struct {
	Title  string `json:"title"`
	TaskID string `json:"task_id"`
	Btns   []*Btn `json:"btn"`

	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Btn struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Color  string `json:"color,omitempty"`
	IsBold bool   `json:"bool,omitempty"`
}

// Todo
func NewMsgInteractiveCardFromPayload(payload *models.Payload) *Msg {
	return &Msg{}
}
