package templates

import (
	_ "embed"
)

//go:embed default.tmpl
var DefaultTmpl string

//go:embed default.zh.tmpl
var DefaultTmplZH string

//go:embed weixin.tmpl
var DefaultTmplWeixin string

//go:embed weixin.zh.tmpl
var DefaultTmplWeixinZH string

//go:embed weixinapp.tmpl
var DefaultTmplWeixinapp string

//go:embed weixinapp.zh.tmpl
var DefaultTmplWeixinappZH string

//go:embed dingtalk.tmpl
var DefaultTmplDingTalk string

//go:embed dingtalk.zh.tmpl
var DefaultTmplDingTalkZH string

//go:embed feishu.tmpl
var DefaultTmplFeishu string

//go:embed feishu.zh.tmpl
var DefaultTmplFeishuZH string

//go:embed slack.tmpl
var DefaultTmplSlack string

//go:embed slack.zh.tmpl
var DefaultTmplSlackZH string

//go:embed discord-webhook.tmpl
var DefaultTmplDiscordWebhook string

var DefaultTmplByLang = map[string]string{
	"en": DefaultTmpl,
	"zh": DefaultTmplZH,
}

// Must define for every supported channel
var ChannelsDefaultTmplMapByLang = map[string]map[string]string{
	"en": {
		"dingtalk":        DefaultTmplDingTalk,
		"feishu":          DefaultTmplFeishu,
		"slack":           DefaultTmplSlack,
		"weixin":          DefaultTmplWeixin,
		"weixinapp":       DefaultTmplWeixinapp,
		"discord-webhook": DefaultTmplDiscordWebhook,
	},
	"zh": {
		"dingtalk":  DefaultTmplDingTalkZH,
		"feishu":    DefaultTmplFeishuZH,
		"slack":     DefaultTmplSlackZH,
		"weixin":    DefaultTmplWeixinZH,
		"weixinapp": DefaultTmplWeixinappZH,
	},
}

func DefaultSupportedLangs() []string {
	res := make([]string, 0)
	for k := range DefaultTmplByLang {
		res = append(res, k)
	}
	return res
}
