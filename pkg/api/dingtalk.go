package api

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/dingtalk"
	restful "github.com/emicklei/go-restful/v3"
)

func createDingtalkSender(request *restful.Request) (*dingtalk.Sender, error) {
	accessToken := request.QueryParameter("access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("not access_token found for dingtalk channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || dingtalk.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for dingtalk")
	}

	return dingtalk.NewSender(accessToken, msgType), nil
}
