package weixinapp

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeMPNews] = NewMsgMPNewsFromPayload
}

type MPNews struct {
	Articles []*Article `json:"articles"` // 图文消息，一个图文消息支持1到8条图文
}

type ArticleForMPNews struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	Content          string `json:"content,omitempty"`
	Digest           string `json:"digest,omitempty"`
}

// Todo
func NewMsgMPNewsFromPayload(payload *models.Payload) *Msg {
	return &Msg{}
}
