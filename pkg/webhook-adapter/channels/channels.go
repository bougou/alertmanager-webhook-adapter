package channels

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/dingtalk"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/feishu"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/slack"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/weixin"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/channels/weixinapp"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func NewDingtalkSender(token string, msgType string) models.Sender {
	return dingtalk.NewSender(token, msgType)
}

func NewFeishuSender(token string, msgType string) models.Sender {
	return feishu.NewSender(token, msgType)
}

func NewWeixinSender(token string, msgType string) models.Sender {
	return weixin.NewSender(token, msgType)
}

func NewWeixinAppSender(corpID string, agentID int, agentSecret string, msgType string, toUser string, toParty string, toTag string) models.Sender {
	return weixinapp.NewSender(corpID, agentID, agentSecret, msgType, toUser, toParty, toTag)
}

func NewSlackSender(token string, channel string, msgType string) models.Sender {
	return slack.NewSender(token, channel, msgType)
}
