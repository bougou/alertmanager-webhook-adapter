package senders

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/feishu"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeFeishu = "feishu"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeFeishu, createFeishuSender)
}

func createFeishuSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for feishu channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}

	var sender models.Sender = feishu.NewSender(token, msgType)
	return sender, nil
}
