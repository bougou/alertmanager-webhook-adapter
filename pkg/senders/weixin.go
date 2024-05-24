package senders

import (
	"fmt"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/channels/weixin"
	"github.com/bougou/webhook-adapter/models"

	restful "github.com/emicklei/go-restful/v3"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeWeixin, createWeixinSender)

	RegisterChannelsMsgConverter(ChannelTypeWeixin, weixin.MsgTypeMarkdown, &weixinMarkdownConverter{})
}

const (
	ChannelTypeWeixin string = "weixin"
)

type (
	weixinMarkdownConverter struct{}
)

func createWeixinSender(request *restful.Request) (sender models.Sender, converter MsgConverter, err error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, nil, fmt.Errorf("not token found for weixin channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = weixin.MsgTypeMarkdown
	}

	sender = weixin.NewSender(token, msgType)

	converter, ok := getMsgConverter(ChannelTypeWeixin, msgType)
	if !ok {
		err = ErrNotFoundConverter(ChannelTypeWeixin, msgType)
		return
	}

	return
}

func (c *weixinMarkdownConverter) Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error) {
	payload, err := promMsg.ToPayload(ChannelTypeWeixin, raw)
	if err != nil {
		return nil, fmt.Errorf("create msg payload failed, %v", err)
	}

	msgType := weixin.MsgTypeMarkdown
	payload2MsgFn, exists := weixin.Payload2MsgFnMap[msgType]
	if !exists {
		return nil, ErrNotFoundPayload2MsgFn(ChannelTypeWeixin, msgType)
	}
	msg := payload2MsgFn(payload)

	return msg, nil
}
