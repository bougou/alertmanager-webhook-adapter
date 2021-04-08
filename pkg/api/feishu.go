package api

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/feishu"
	restful "github.com/emicklei/go-restful/v3"
)

func createFeishuSender(request *restful.Request) (*feishu.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for feishu channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || feishu.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for feishu")
	}

	return feishu.NewSender(token, msgType), nil
}
