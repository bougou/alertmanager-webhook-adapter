package models

import (
	"strings"
	"testing"
	"time"
)

func testAlert(labels KV) *AlertmanagerWebhookMessage {
	return &AlertmanagerWebhookMessage{
		Status:      "firing",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{
				Status: "firing",
				Labels: labels,
				Annotations: KV{
					"description": "cpu usage is high than 20% for 5 minutes",
				},
				StartsAt: time.Date(2021, 3, 30, 20, 17, 50, 0, time.Local),
			},
		},
	}
}

func TestAlertInstanceFallsBackToInstanceLabel(t *testing.T) {
	msg := testAlert(KV{
		"alertname": "CPUHigh",
		"severity":  "warning",
		"instance":  "10.30.1.160",
	})

	channels := []string{"dingtalk", "feishu", "slack", "weixin", "weixinapp", "discord-webhook", "default"}
	for _, ch := range channels {
		out, err := msg.RenderTmpl(ch, "prom.markdown")
		if err != nil {
			t.Fatalf("channel %s: %v", ch, err)
		}
		if !strings.Contains(out, "10.30.1.160") {
			t.Errorf("channel %s: expected instance 10.30.1.160 in output, got:\n%s", ch, out)
		}
	}
}

func TestAlertInstancePrefersAlertinstanceLabel(t *testing.T) {
	msg := testAlert(KV{
		"alertname":     "CPUHigh",
		"severity":      "warning",
		"alertinstance": "ns/pod",
		"instance":      "10.30.1.160",
	})

	out, err := msg.RenderTmpl("dingtalk", "prom.markdown")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "ns/pod") {
		t.Errorf("expected alertinstance ns/pod in output, got:\n%s", out)
	}
	if strings.Contains(out, "10.30.1.160") {
		t.Errorf("instance should not appear when alertinstance is set, got:\n%s", out)
	}
}
