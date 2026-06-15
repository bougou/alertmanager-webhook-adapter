package dingtalk

import "github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"

func init() {
	Payload2MsgFnMap[MsgTypeLink] = NewMsgLinkFromPayload
}

type Link struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicURL     string `json:"picUrl"`
	MessageURL string `json:"messageUrl"`
}

func NewLink(title string, text string, messageURL string) *Link {
	return &Link{
		Text:       text,
		Title:      title,
		MessageURL: messageURL,
	}
}

func (link *Link) WithPicURL(picURL string) *Link {
	link.PicURL = picURL
	return link
}

func NewMsgLink(link *Link) *Msg {
	return &Msg{
		MsgType: MsgTypeLink,
		Link:    link,
	}
}

func NewMsgLinkFromPayload(payload *models.Payload) *Msg {
	link := NewLink(payload.Title, payload.Text, "")
	return NewMsgLink(link)
}
