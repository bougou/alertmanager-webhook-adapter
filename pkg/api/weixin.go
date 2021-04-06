package api

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/dingtalk"
	"github.com/bougou/webhook-adapter/channels/weixin"
	restful "github.com/emicklei/go-restful/v3"
)

func createWeixinSender(request *restful.Request) (*weixin.Sender, error) {
	key := request.QueryParameter("key")
	if key == "" {
		return nil, fmt.Errorf("not access_token found for dingtalk channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || dingtalk.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for dingtalk")

	}

	return weixin.NewSender(key, msgType), nil
}
