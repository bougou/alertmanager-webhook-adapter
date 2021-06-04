package senders

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/weixin"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeWeixin = "weixin"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeWeixin, createWeixinSender)
}

func createWeixinSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for weixin channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || weixin.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for weixin")

	}

	var sender models.Sender = weixin.NewSender(token, msgType)
	return sender, nil
}
