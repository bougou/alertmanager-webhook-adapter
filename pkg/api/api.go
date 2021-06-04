package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	promModels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/senders"
	restful "github.com/emicklei/go-restful/v3"
)

type Controller struct {
	signature string
}

func NewController(signature string) *Controller {
	return &Controller{
		signature: signature,
	}
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
	promMsg.SetMessageAt().SetSignature(c.signature)

	channelType := request.QueryParameter("channel_type")
	if channelType == "" {
		response.WriteHeaderAndJson(http.StatusBadRequest, "not channel_type found", restful.MIME_JSON)
		return
	}

	senderCreator, exists := senders.ChannelsSenderCreatorMap[channelType]
	if !exists {
		response.WriteHeaderAndJson(http.StatusBadRequest, "not supported channel_type", restful.MIME_JSON)
		return
	}

	sender, err := senderCreator(request)
	if err != nil {
		response.WriteHeaderAndJson(http.StatusInternalServerError, fmt.Sprintf("create sender failed, %v", err), restful.MIME_JSON)
		return
	}

	payload := promMsg.ToPayload(channelType, raw)

	if err := sender.Send(payload); err != nil {
		response.WriteHeaderAndJson(http.StatusInternalServerError, fmt.Sprintf("sender send failed, %v", err), restful.MIME_JSON)
		return
	}

}
