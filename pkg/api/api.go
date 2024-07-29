package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	promModels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/senders"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/kr/pretty"
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
	c.logf("Got request : %s\n", request.Request.URL.String())

	raw, err := io.ReadAll(request.Request.Body)
	if err != nil {
		errmsg := fmt.Sprintf("Err: read request body failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	promMsg := &promModels.AlertmanagerWebhookMessage{}
	if err := json.Unmarshal(raw, promMsg); err != nil {
		errmsg := fmt.Sprintf("Err: unmarshal body failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}
	promMsg.SetMessageAt().SetSignature(c.signature)

	channelType := request.QueryParameter("channel_type")
	if channelType == "" {
		errmsg := "Err: no channel_type found"
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	senderCreator, exists := senders.ChannelsSenderCreatorMap[channelType]
	if !exists {
		errmsg := fmt.Sprintf("Err: not supported channel_type of (%s)", channelType)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	sender, converter, err := senderCreator(request)
	if err != nil {
		errmsg := fmt.Sprintf("Err: create sender failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusBadRequest, errmsg, restful.MIME_JSON)
		return
	}

	msg, err := converter.Convert(raw, promMsg)
	if err != nil {
		errmsg := fmt.Sprintf("Err: convert the prom msg failed, err: %s", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusInternalServerError, errmsg, restful.MIME_JSON)
		return
	}
	if c.debug {
		pretty.Println(payload)

		fmt.Println(">>> Payload Markdown")
		fmt.Print(payload.Markdown)
	}

	if err := sender.SendMsg(msg); err != nil {
		errmsg := fmt.Sprintf("Err: sender send failed, %v", err)
		c.log(errmsg)
		response.WriteHeaderAndJson(http.StatusInternalServerError, errmsg, restful.MIME_JSON)
		return
	}

	c.logf("Send succeed: %s\n", request.Request.URL.String())
	response.WriteHeader(http.StatusNoContent)
}
