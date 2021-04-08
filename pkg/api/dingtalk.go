package api

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/dingtalk"
	restful "github.com/emicklei/go-restful/v3"
)

func createDingtalkSender(request *restful.Request) (*dingtalk.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for dingtalk channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || dingtalk.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for dingtalk")
	}

	return dingtalk.NewSender(token, msgType), nil
}
