package api

import (
	"fmt"

	"github.com/bougou/webhook-adapter/channels/weixin"
	restful "github.com/emicklei/go-restful/v3"
)

func createWeixinSender(request *restful.Request) (*weixin.Sender, error) {
	token := request.QueryParameter("token")
	if token == "" {
		return nil, fmt.Errorf("not token found for weixin channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || weixin.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for weixin")

	}

	return weixin.NewSender(token, msgType), nil
}
