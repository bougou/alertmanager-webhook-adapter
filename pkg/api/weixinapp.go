package api

import (
	"fmt"
	"strconv"

	"github.com/bougou/webhook-adapter/channels/weixinapp"
	restful "github.com/emicklei/go-restful/v3"
)

func createWeixinappSender(request *restful.Request) (*weixinapp.Sender, error) {
	corpID := request.QueryParameter("corp_id")
	if corpID == "" {
		return nil, fmt.Errorf("not core_id found for weixin channel")
	}

	agentID := request.QueryParameter("agent_id")
	if agentID == "" {
		return nil, fmt.Errorf("not agent_id found for weixin channel")
	}

	aID, err := strconv.Atoi(agentID)
	if err != nil {
		return nil, fmt.Errorf("agent_id must be integer")
	}

	agentSecret := request.QueryParameter("agent_secret")
	if agentSecret == "" {
		return nil, fmt.Errorf("not agent_secret found for weixin channel")
	}

	msgType := request.QueryParameter("msg_type")
	if !(msgType == "" || weixinapp.ValidMsgtype(msgType)) {
		return nil, fmt.Errorf("not supported msgtype for weixin")
	}

	return weixinapp.NewSender(corpID, aID, agentSecret, msgType), nil
}
