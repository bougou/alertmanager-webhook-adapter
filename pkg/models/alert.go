package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bougou/webhook-adapter/models"
)

// ref: https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
// The Alertmanager will send HTTP POST requests in the following JSON format to the configured endpoint:
// {
//   "version": "4",
//   "groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
//   "truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
//   "status": "<resolved|firing>",
//   "receiver": <string>,
//   "groupLabels": <object>,
//   "commonLabels": <object>,
//   "commonAnnotations": <object>,
//   "externalURL": <string>,           // backlink to the Alertmanager.
//   "alerts": [
//     {
//       "status": "<resolved|firing>",
//       "labels": <object>,
//       "annotations": <object>,
//       "startsAt": "<rfc3339>",
//       "endsAt": "<rfc3339>",
//       "generatorURL": <string>       // identifies the entity that caused the alert
//     },
//     ...
//   ]
// }

type AlertmanagerWebhookMessage struct {
	Version           string           `json:"version"`
	GroupKey          *json.RawMessage `json:"groupKey"`
	TruncatedAlerts   int              `json:"truncatedAlerts"`
	Status            string           `json:"status"`
	Receiver          string           `json:"receiver"`
	GroupLabels       KV               `json:"groupLabels"`
	CommonLabels      KV               `json:"commonLabels"`
	CommonAnnotations KV               `json:"commonAnnotations"`
	ExternalURL       string           `json:"externalURL"`
	Alerts            Alerts           `json:"alerts"`

	// extra fields added by us
	MessageAt time.Time `json:"messageAt"` // the time the webhook message was received
	Signature string    `json:"signature"` // 签名，如发送短信时出现在内容最前面【】
}

type Alerts []Alert

type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
}

// Firing Alerts the subset of alerts that are firing.
func (as Alerts) Firing() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == "firing" {
			res = append(res, a)
		}
	}
	return res
}

// Resolved returns the subset of alerts that are resolved.
func (as Alerts) Resolved() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == "resolved" {
			res = append(res, a)
		}
	}
	return res
}

func (m *AlertmanagerWebhookMessage) SetMessageAt() *AlertmanagerWebhookMessage {
	m.MessageAt = time.Now()
	return m
}

func (m *AlertmanagerWebhookMessage) SetSignature(s string) *AlertmanagerWebhookMessage {
	m.Signature = s
	return m
}

func (m *AlertmanagerWebhookMessage) RenderTmpl(tmplName string) (string, error) {
	tmpl, err := promTemplate.Clone()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, m); err != nil {
		return "<<<< template error >>>>", err
	}

	return buf.String(), nil
}

func (m *AlertmanagerWebhookMessage) ToPayload(raw []byte) *models.Payload {
	msg := &models.Payload{Raw: string(raw)}

	title, err := m.RenderTmpl("prom.title")
	if err != nil {
		fmt.Println(err)
	}
	msg.Title = title

	text, err := m.RenderTmpl("prom.text")
	if err != nil {
		fmt.Println(err)
	}
	msg.Text = text

	markdown, err := m.RenderTmpl("prom.markdown")
	if err != nil {
		fmt.Println(err)
	}
	msg.Markdown = markdown

	return msg
}
