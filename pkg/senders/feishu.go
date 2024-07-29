package senders

import (
	"fmt"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/channels/feishu"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeFeishu = "feishu"
)

type (
	feishuMarkdownConverter struct{}
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeFeishu, createFeishuSender)

	RegisterChannelsMsgConverter(ChannelTypeFeishu, feishu.MsgTypeMarkdown, &feishuMarkdownConverter{})
}

func createFeishuSender(request *restful.Request) (sender models.Sender, converter MsgConverter, err error) {
	token := request.QueryParameter("token")
	if token == "" {
		err = fmt.Errorf("not token found for feishu channel")
		return
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	sender = feishu.NewSender(token, msgType)

	converter, ok := getMsgConverter(ChannelTypeFeishu, msgType)
	if !ok {
		err = ErrNotFoundConverter(ChannelTypeFeishu, msgType)
		return
	}

	return
}

func (c *feishuMarkdownConverter) Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error) {
	payload, err := promMsg.ToPayload(ChannelTypeFeishu, raw)
	if err != nil {
		return nil, fmt.Errorf("create msg payload failed, %v", err)
	}

	msgType := feishu.MsgTypeMarkdown
	payload2MsgFn, exists := feishu.Payload2MsgFnMap[msgType]
	if !exists {
		return nil, ErrNotFoundPayload2MsgFn(ChannelTypeFeishu, msgType)
	}
	msg := payload2MsgFn(payload)

	return msg, nil
}
