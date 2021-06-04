package templates

import (
	_ "embed"
)

//go:embed default.tmpl
var DefaultTmpl string

//go:embed weixin.tmpl
var DefaultTmplWeixin string

//go:embed weixinapp.tmpl
var DefaultTmplWeixinapp string

//go:embed dingtalk.tmpl
var DefaultTmplDingTalk string

//go:embed feishu.tmpl
var DefaultTmplFeishu string

// Must define for every supported channel
var ChannelsDefaultTmplMap = map[string]string{
	"weixin":    DefaultTmplWeixin,
	"weixinapp": DefaultTmplWeixinapp,
	"dingtalk":  DefaultTmplDingTalk,
	"feishu":    DefaultTmplFeishu,
}
