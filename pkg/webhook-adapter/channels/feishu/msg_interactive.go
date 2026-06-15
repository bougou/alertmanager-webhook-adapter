package feishu

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/feishu/card"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeInteractive] = NewMsgInteractiveFromPayload
}

type Card struct {
	Config       *CardConfig       `json:"config,omitempty"`
	Header       *CardHeader       `json:"header,omitempty"`
	CardLink     *card.MultiURL    `json:"card_link,omitempty"`
	Elements     []card.CardModule `json:"elements"` // 最多可堆叠 50 个模块
	I18NElements *I18NElements     `json:"i18n_elements,omitempty"`
}

// CardConfig 卡片配置
type CardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode,omitempty"` // 2021/03/22 之后，此字段废弃，所有卡片均升级为自适应屏幕宽度的宽版卡片
	EnableForward  bool `json:"enable_forward"`             // 是否允许卡片被转发，默认 false
}

type CardHeader struct {
	Title    *card.Text `json:"title"`              // 卡片标题内容, text 对象（仅支持 "plain_text")
	Template string     `json:"template,omitempty"` // 控制标题背景颜色, https://open.feishu.cn/document/ukTMukTMukTM/ukTNwUjL5UDM14SO1ATN
}

type I18NElements struct {
	ZHCN []card.CardModule `json:"zh_cn"`
	ENUS []card.CardModule `json:"en_us"`
	JAJP []card.CardModule `json:"jn_jp"`
}

func NewMsgInteractive(card *Card) *Msg {
	return &Msg{
		MsgType: MsgTypeInteractive,
		Card:    card,
	}
}

func NewMsgInteractiveFromPayload(payload *models.Payload) *Msg {
	// Todo, construct card from payload
	card := &Card{}
	return NewMsgInteractive(card)
}
