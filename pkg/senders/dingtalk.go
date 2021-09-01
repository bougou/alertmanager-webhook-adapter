package senders

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/dingtalk"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	ChannelTypeDingtalk = "dingtalk"
)

func init() {
	RegisterChannelsSenderCreator(ChannelTypeDingtalk, createDingtalkSender)
}

func createDingtalkSender(request *restful.Request) (models.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for dingtalk channel")
	}

	msgType := request.QueryParameter("msg_type")
	if msgType == "" {
		msgType = "markdown"
	}
	if !(dingtalk.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for dingtalk")
	}

	var sender models.Sender = dingtalk.NewSender(token, msgType)
	return sender, nil
}
