package models

import (
	"bytes"
	"encoding/json"
	"errors"
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

// AlertManagerWebhookMessage holds data that alertmanager passed to webhook server.
// https://pkg.go.dev/github.com/prometheus/alertmanager/template#Data already defines
// a struct type, but here we re-defined it, cause we will fill extra fields into it.
type AlertmanagerWebhookMessage struct {
	Version         string           `json:"version"`
	GroupKey        *json.RawMessage `json:"groupKey"`
	TruncatedAlerts int              `json:"truncatedAlerts"`

	Status            string `json:"status"`
	Receiver          string `json:"receiver"`
	Alerts            Alerts `json:"alerts"`
	GroupLabels       KV     `json:"groupLabels"`
	CommonLabels      KV     `json:"commonLabels"`
	CommonAnnotations KV     `json:"commonAnnotations"`
	ExternalURL       string `json:"externalURL"`

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

// Firing returns the subset of alerts that are firing.
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

func (alert *Alert) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, v := range m {
		switch k {
		case "startsAt":
			_t := v.(string)
			t, err := parseTimeFromStr(_t)
			if err != nil {
				return err
			}
			alert.StartsAt = t

		case "endsAt":
			_t := v.(string)
			t, err := parseTimeFromStr(_t)
			if err != nil {
				return err
			}
			alert.EndsAt = t

		case "status":
			alert.Status = v.(string)

		case "generatorURL":
			alert.GeneratorURL = v.(string)

		case "labels":
			s, err := json.Marshal(v)
			if err != nil {
				return err
			}
			kv := KV{}
			if err := json.Unmarshal(s, &kv); err != nil {
				return err
			}
			alert.Labels = kv

		case "annotations":
			s, err := json.Marshal(v)
			if err != nil {
				return err
			}
			kv := KV{}
			if err := json.Unmarshal(s, &kv); err != nil {
				return err
			}
			alert.Annotations = kv
		}
	}

	return nil
}

func parseTimeFromStr(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return t, err
	}

	return t.In(time.Local), nil
}

func (m *AlertmanagerWebhookMessage) SetMessageAt() *AlertmanagerWebhookMessage {
	m.MessageAt = time.Now()
	return m
}

func (m *AlertmanagerWebhookMessage) SetSignature(s string) *AlertmanagerWebhookMessage {
	m.Signature = s
	return m
}

func (m *AlertmanagerWebhookMessage) RenderTmpl(channel string, tmplName string) (string, error) {
	var safetmpl *safeTemplate

	if t, exists := promMsgTemplatesMap[channel]; !exists {
		safetmpl = promMsgTemplateDefault
	} else {
		safetmpl = t
	}

	tmpl, err := safetmpl.Clone()
	if err != nil {
		msg := fmt.Sprintf("Clone template failed, err: %s", err)
		return "", errors.New(msg)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, m); err != nil {
		msg := fmt.Sprintf("ExecuteTemplate failed, err: %s", err)
		return "<<<< template error >>>>", errors.New(msg)
	}

	return buf.String(), nil
}

func (m *AlertmanagerWebhookMessage) ToPayload(channel string, raw []byte) (*models.Payload, error) {
	payload := &models.Payload{Raw: string(raw)}

	title, err := m.RenderTmpl(channel, "prom.title")
	if err != nil {
		return nil, err
	}
	payload.Title = title

	text, err := m.RenderTmpl(channel, "prom.text")
	if err != nil {
		return nil, err
	}
	payload.Text = text

	markdown, err := m.RenderTmpl(channel, "prom.markdown")
	if err != nil {
		return nil, err
	}
	payload.Markdown = markdown

	return payload, nil
}
