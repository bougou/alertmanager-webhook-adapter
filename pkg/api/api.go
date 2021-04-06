package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bougou/alertmanager-webhook-adapter/cmd/alertmanager-webhook-adapter/app/options"
	promModels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

type Controller struct {
	appOptions *options.AppOptions
}

func NewController(o *options.AppOptions) *Controller {
	return &Controller{o}
}

func (c *Controller) Install(container *restful.Container) {

	ws := new(restful.WebService)
	ws.Path("/webhook/send")

	ws.Route(
		ws.POST("/").To(c.send),
	)

	container.Add(ws)
}

func (c *Controller) send(request *restful.Request, response *restful.Response) {

	raw, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		response.WriteHeaderAndJson(http.StatusBadRequest, "read request body failed", restful.MIME_JSON)
		return
	}

	promMsg := &promModels.AlertmanagerWebhookMessage{}
	if err := json.Unmarshal(raw, promMsg); err != nil {
		response.WriteHeaderAndJson(http.StatusBadRequest, "unmarshal body failed", restful.MIME_JSON)
		return
	}
	promMsg.SetMessageAt().SetSignature(c.appOptions.Signature)

	payload := promMsg.ToPayload(raw)

	channelType := request.QueryParameter("channel_type")
	if channelType == "" {
		response.WriteHeaderAndJson(http.StatusBadRequest, "not channel_type found", restful.MIME_JSON)
		return
	}

	var sender models.Sender

	switch channelType {
	case "weixin":
		sender, err = createWeixinSender(request)
	case "dingtalk":
		sender, err = createDingtalkSender(request)
	case "feishu":
		sender, err = createFeishuSender(request)
	default:
		response.WriteHeaderAndJson(http.StatusBadRequest, "not supported channel_type", restful.MIME_JSON)
		return
	}

	if err != nil {
		response.WriteHeaderAndJson(http.StatusInternalServerError, fmt.Sprintf("create sender failed, %v", err), restful.MIME_JSON)
		return
	}

	if err := sender.Send(payload); err != nil {
		response.WriteHeaderAndJson(http.StatusInternalServerError, fmt.Sprintf("sender send failed, %v", err), restful.MIME_JSON)
		return
	}

}
