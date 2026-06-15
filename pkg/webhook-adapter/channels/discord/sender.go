package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	"github.com/bwmarrin/discordgo"
)

// https://discord.com/api/webhooks/XXX/YYY
//
//	XXX is the webhookID
//	YYY is the webhookAuthToken
//
// https://discord.com/api/webhooks/1257662116761174016/YDmFPVBP4XJ7UNveeowbcv78PySW5GRI-yStxBXgo2El4y64snfXaHZDCsKcXvNK8f8x
type WebhookSender struct {
	ID    string // webhook ID
	Token string // webhook auth token
}

var _ models.Sender = (*WebhookSender)(nil)

func NewWebhookSender(id, token string) *WebhookSender {
	return &WebhookSender{
		ID:    id,
		Token: token,
	}
}

func (ws *WebhookSender) Addr() string {
	return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", ws.ID, ws.Token)
}

func (ws *WebhookSender) Send(payload *models.Payload) error {
	msgWebhook := NewMsgWebhookFromPayload(payload)
	return ws.SendMsg(msgWebhook)
}

func (ws *WebhookSender) SendMsg(msgSource interface{}) error {
	return ws.SendMsgT("not required", msgSource)
}

func (ws *WebhookSender) SendMsgT(msgType string, msgSource interface{}) error {
	msg, ok := msgSource.(*MsgWebhook)
	if !ok {
		return fmt.Errorf("passed msgSource is not type *MsgWebhook")
	}

	// when sending msg to webhook, empty session needed
	s, err := discordgo.New("")
	if err != nil {
		return fmt.Errorf("new discordgo session failed, err: %s", err)
	}
	_, err = s.WebhookExecute(ws.ID, ws.Token, false, &msg.WebhookParams)
	if err != nil {
		return fmt.Errorf("webhook execute failed, err: %s", err)
	}

	return nil
}

func (ws *WebhookSender) SendByDirectHTTP(msg *MsgWebhook) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(msgBytes))

	req, err := http.NewRequest("POST", ws.Addr(), bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("failed to construct request, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send msg error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read res body failed, err: %s", err)
	}

	if len(resBody) == 0 {
		return fmt.Errorf("invalid response")
	}

	type DiscordWebhookResponse struct {
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	}

	dwRes := &DiscordWebhookResponse{}
	if err := json.Unmarshal(resBody, dwRes); err != nil {
		return fmt.Errorf("marshal to DiscordWebhookResponse failed, err: %s", err)
	}

	if dwRes.Code != "0" {
		return fmt.Errorf("found err in response, request url: (%s), response: %s", ws.Addr(), string(resBody))
	}

	return nil
}
