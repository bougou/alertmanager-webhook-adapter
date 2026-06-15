package weixin

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeTemplateCard] = NewMsgTemplateCardFromPayload
}

type TemplateCard struct {
	CardType              string                          `json:"card_type,omitempty"`
	Source                TemplateCardSource              `json:"source,omitempty"`
	MainTitle             TemplateCardMainTitle           `json:"main_title,omitempty"`
	EmphasisContent       TemplateCardEmphasisContent     `json:"emphasis_content,omitempty"`
	QuoteArea             TemplateCardQuoteArea           `json:"quote_area,omitempty"`
	SubTitleText          string                          `json:"sub_title_text,omitempty"`
	HorizontalContentList []TemplateCardHorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              []TemplateCardJump              `json:"jump_list,omitempty"`
	CardAction            TemplateCardCardAction          `json:"card_action,omitempty"`
}

type TemplateCardSource struct {
	IconURL   string `json:"icon_url,omitempty"`
	Desc      string `json:"desc,omitempty"`
	DescColor int    `json:"desc_color,omitempty"`
}

type TemplateCardMainTitle struct {
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type TemplateCardEmphasisContent struct {
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type TemplateCardQuoteArea struct {
	Type      int    `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
	Appid     string `json:"appid,omitempty"`
	Pagepath  string `json:"pagepath,omitempty"`
	Title     string `json:"title,omitempty"`
	QuoteText string `json:"quote_text,omitempty"`
}

type TemplateCardHorizontalContent struct {
	Keyname string `json:"keyname,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    int    `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	MediaID string `json:"media_id,omitempty"`
}

type TemplateCardJump struct {
	Type     int    `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
	Title    string `json:"title,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

type TemplateCardCardAction struct {
	Type     int    `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

func NewMsgTemplateCard(templateCard *TemplateCard) *Msg {
	return &Msg{
		MsgType:      MsgTypeTemplateCard,
		TemplateCard: templateCard,
	}
}

func NewMsgTemplateCardFromPayload(payload *models.Payload) *Msg {
	var templateCard TemplateCard
	return NewMsgTemplateCard(&templateCard)
}

func (b *WeixinGroupBot) SendTemplateCard(templateCard *TemplateCard) error {
	msg := NewMsgTemplateCard(templateCard)
	return b.Send(msg)
}
