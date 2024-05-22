package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	promModels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/senders"
	restful "github.com/emicklei/go-restful/v3"
)

type Controller struct {
	signature string
	debug     bool
}

func NewController(signature string) *Controller {
	return &Controller{
		signature: signature,
	}
}

func (c *Controller) WithDebug(debug bool) *Controller {
	if debug {
		fmt.Println("debug mode enabled")
	}
	c.debug = debug
	return c
}

func (c *Controller) Install(container *restful.Container) {

	ws := new(restful.WebService)
	ws.Path("/webhook/send")

	ws.Route(
		ws.POST("/").To(c.send),
	)

	container.Add(ws)
}

func (c *Controller) logf(format string, a ...any) error {
	if c.debug {
		_, err := fmt.Printf(format, a...)
		return err
	}

	return nil
}

func (c *Controller) log(a ...any) error {
	if c.debug {
		_, err := fmt.Println(a...)
		return err
	}

	return nil
}

func (c *Controller) send(request *restful.Request, response *restful.Response) {
	c.logf("Got: %s\n", request.Request.URL.String())

	raw, err := io.ReadAll(request.Request.Body)
	if err != nil {
		errmsg := fmt.Sprintf("read request body failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	promMsg := &promModels.AlertmanagerWebhookMessage{}
	if err := json.Unmarshal(raw, promMsg); err != nil {
		errmsg := fmt.Sprintf("unmarshal body failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}
	promMsg.SetMessageAt().SetSignature(c.signature)

	channelType := request.QueryParameter("channel_type")
	if channelType == "" {
		errmsg := "no channel_type found"
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	senderCreator, exists := senders.ChannelsSenderCreatorMap[channelType]
	if !exists {
		errmsg := "not supported channel_type"
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	sender, err := senderCreator(request)
	if err != nil {
		errmsg := fmt.Sprintf("create sender failed, %v", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	payload, err := promMsg.ToPayload(channelType, raw)
	if err != nil {
		errmsg := fmt.Sprintf("create msg payload failed, %v", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusInternalServerError, errmsg, restful.MIME_JSON)
		return
	}

	if err := sender.Send(payload); err != nil {
		errmsg := fmt.Sprintf("sender send failed, %v", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusInternalServerError, errmsg, restful.MIME_JSON)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
