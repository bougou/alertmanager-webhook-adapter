package senders

import (
	"fmt"
	"strconv"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/channels/weixinapp"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeWeixinApp = "weixinapp"
)

type (
	weixinappMarkdownConverter struct{}
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeWeixinApp, createWeixinappSender)

	RegisterChannelsMsgConverter(ChannelTypeWeixinApp, weixinapp.MsgTypeMarkdown, &weixinappMarkdownConverter{})
}

func createWeixinappSender(request *restful.Request) (sender models.Sender, converter MsgConverter, err error) {
	corpID := request.QueryParameter("corp_id")
	if corpID == "" {
		err = fmt.Errorf("not core_id found for weixin channel")
		return
	}

	agentID := request.QueryParameter("agent_id")
	if agentID == "" {
		err = fmt.Errorf("not agent_id found for weixin channel")
		return
	}

	aID, err := strconv.Atoi(agentID)
	if err != nil {
		err = fmt.Errorf("agent_id must be integer")
		return
	}

	agentSecret := request.QueryParameter("agent_secret")
	if agentSecret == "" {
		err = fmt.Errorf("not agent_secret found for weixin channel")
		return
	}

	toUser := request.QueryParameter("to_user")
	toParty := request.QueryParameter("to_party")
	toTag := request.QueryParameter("to_tag")

	if toUser == "" && toParty == "" && toTag == "" {
		err = fmt.Errorf("must specify one of to_user,to_party,to_tag")
		return
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	sender = weixinapp.NewSender(corpID, aID, agentSecret, msgType, toUser, toParty, toTag)

	converter, ok := getMsgConverter(ChannelTypeWeixinApp, msgType)
	if !ok {
		err = ErrNotFoundConverter(ChannelTypeWeixinApp, msgType)
		return
	}

	return
}

func (c *weixinappMarkdownConverter) Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error) {
	payload, err := promMsg.ToPayload(ChannelTypeWeixinApp, raw)
	if err != nil {
		return nil, fmt.Errorf("create msg payload failed, %v", err)
	}

	msgType := weixinapp.MsgTypeMarkdown
	payload2MsgFn, exists := weixinapp.Payload2MsgFnMap[msgType]
	if !exists {
		return nil, ErrNotFoundPayload2MsgFn(ChannelTypeWeixinApp, msgType)
	}
	msg := payload2MsgFn(payload)

	return msg, nil
}
