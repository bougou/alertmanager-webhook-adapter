package senders

import (
	"fmt"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/channels/dingtalk"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeDingtalk = "dingtalk"
)

type (
	dingtalkMarkdownConverter struct{}
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeDingtalk, createDingtalkSender)

	RegisterChannelsMsgConverter(ChannelTypeDingtalk, dingtalk.MsgTypeMarkdown, &dingtalkMarkdownConverter{})
}

func createDingtalkSender(request *restful.Request) (sender models.Sender, converter MsgConverter, err error) {
	token := request.QueryParameter("token")
	if token == "" {
		err = fmt.Errorf("not token found for dingtalk channel")
		return
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	sender = dingtalk.NewSender(token, msgType)

	converter, ok := getMsgConverter(ChannelTypeDingtalk, msgType)
	if !ok {
		err = ErrNotFoundConverter(ChannelTypeDingtalk, msgType)
		return
	}

	return
}

func (c *dingtalkMarkdownConverter) Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error) {
	payload, err := promMsg.ToPayload(ChannelTypeDingtalk, raw)
	if err != nil {
		return nil, fmt.Errorf("create msg payload failed, %v", err)
	}

	msgType := dingtalk.MsgTypeMarkdown
	payload2MsgFn, exists := dingtalk.Payload2MsgFnMap[msgType]
	if !exists {
		return nil, ErrNotFoundPayload2MsgFn(ChannelTypeDingtalk, msgType)
	}
	msg := payload2MsgFn(payload)

	return msg, nil
}
