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

func TestHighestSeverityPrefersFiringCritical(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Alerts: Alerts{
			{Status: "resolved", Labels: KV{"severity": "warning"}},
			{Status: "firing", Labels: KV{"severity": "warning"}},
			{Status: "firing", Labels: KV{"severity": "critical"}},
		},
	}
	if got := msg.HighestSeverity(); got != "critical" {
		t.Fatalf("HighestSeverity() = %q, want critical", got)
	}
}

func TestHighestSeverityFallsBackToResolved(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Alerts: Alerts{
			{Status: "resolved", Labels: KV{"severity": "info"}},
			{Status: "resolved", Labels: KV{"severity": "warning"}},
		},
	}
	if got := msg.HighestSeverity(); got != "warning" {
		t.Fatalf("HighestSeverity() = %q, want warning", got)
	}
}

func TestToPayloadCopiesAlertmanagerStatusAndSeverity(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Status:      "firing",
		Signature:   "ops",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{Status: "resolved", Labels: KV{"severity": "warning", "alertname": "CPUHigh"}},
			{Status: "firing", Labels: KV{"severity": "critical", "alertname": "CPUHigh"}},
		},
	}
	payload, err := msg.ToPayload("feishu", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if payload.Status != "firing" {
		t.Fatalf("payload.Status = %q, want firing", payload.Status)
	}
	if payload.Severity != "critical" {
		t.Fatalf("payload.Severity = %q, want critical", payload.Severity)
	}
}

func TestFeishuTitleUsesHighestSeverity(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{Status: "firing", Labels: KV{"severity": "warning", "alertname": "CPUHigh"}},
			{Status: "firing", Labels: KV{"severity": "critical", "alertname": "CPUHigh"}},
		},
	}
	out, err := msg.RenderTmpl("feishu", "prom.title")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "CRITICAL") {
		t.Fatalf("title should use highest severity CRITICAL, got: %s", out)
	}
	if strings.Contains(out, "WARNING") {
		t.Fatalf("title should not keep first-alert WARNING, got: %s", out)
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

func TestFeishuMarkdownKeepsVerticalFields(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		Receiver:    "default",
		ExternalURL: "http://am.example.com",
		MessageAt:   time.Date(2026, 4, 17, 18, 3, 40, 0, time.Local),
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{
				Status: "firing",
				Labels: KV{
					"alertname": "CPUHigh",
					"severity":  "warning",
					"instance":  "10.30.1.160",
					"region":    "region-x",
				},
				Annotations: KV{
					"description": "cpu usage is high than 20% for 5 minutes",
				},
				StartsAt:     time.Date(2021, 3, 30, 20, 17, 50, 0, time.Local),
				GeneratorURL: "http://prometheus.example.com/graph",
			},
		},
	}

	out, err := msg.RenderTmpl("feishu", "prom.markdown")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"**Alert Name**:",
		"**Alert Level**:",
		"**Instance**:",
		"**Region**:",
		"**Start At**:",
		"**End At**:",
		"**Description**:",
		"10.30.1.160",
		"**Summary**",
		"**Detail**",
		"\n ---\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
	if strings.Contains(out, "**Graph**") || strings.Contains(out, "[Alertmanager]") {
		t.Errorf("body should match original field list without extra links, got:\n%s", out)
	}
	if strings.Contains(out, "**Firing") {
		t.Errorf("Firing subtitle should not be bold like Summary/Detail, got:\n%s", out)
	}
	inst := strings.Index(out, "**Instance**:")
	name := strings.Index(out, "**Alert Name**:")
	if inst < 0 || name < 0 || inst > name {
		t.Errorf("Instance should be the first field of each detail item, got:\n%s", out)
	}
}

func renderFeishuZHMarkdown(t *testing.T, msg *AlertmanagerWebhookMessage) string {
	t.Helper()
	return renderFeishuZH(t, msg, "prom.markdown")
}

func TestFeishuZHMarkdownSeparatesResolvedAlerts(t *testing.T) {
	start := time.Date(2021, 3, 30, 20, 17, 50, 0, time.Local)
	end := time.Date(2021, 3, 30, 21, 17, 50, 0, time.Local)
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		Receiver:    "default",
		ExternalURL: "http://am.example.com",
		MessageAt:   time.Date(2026, 4, 17, 18, 3, 40, 0, time.Local),
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{
				Status: "resolved",
				Labels: KV{
					"alertname": "CPUHigh",
					"severity":  "warning",
					"instance":  "10.30.1.161",
				},
				StartsAt: start,
				EndsAt:   end,
			},
			{
				Status: "resolved",
				Labels: KV{
					"alertname": "CPUHigh",
					"severity":  "warning",
					"instance":  "10.30.1.164",
				},
				StartsAt: start,
				EndsAt:   end,
			},
		},
	}

	out := renderFeishuZHMarkdown(t, msg)

	first := strings.Index(out, "- **实例**: <text_tag color='purple'>10.30.1.161</text_tag>")
	second := strings.Index(out, "- **实例**: <text_tag color='purple'>10.30.1.164</text_tag>")
	if first < 0 || second < 0 || first > second {
		t.Fatalf("expected both resolved instances in detail list, got:\n%s", out)
	}
	if !strings.Contains(out[first:second], "\n\n") {
		t.Fatalf("resolved alerts should be separated by a blank line, got:\n%s", out)
	}
}

func TestFeishuZHMarkdownSeparatesSummaryTitleFromList(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		Receiver:    "default",
		ExternalURL: "http://am.example.com",
		MessageAt:   time.Date(2026, 4, 17, 18, 3, 40, 0, time.Local),
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{
				Status: "firing",
				Labels: KV{
					"alertname": "CPUHigh",
					"severity":  "warning",
					"instance":  "10.30.1.160",
				},
			},
			{
				Status: "resolved",
				Labels: KV{
					"alertname": "CPUHigh",
					"severity":  "warning",
					"instance":  "10.30.1.161",
				},
			},
		},
	}

	out := renderFeishuZHMarkdown(t, msg)
	for _, title := range []string{
		"<font color='red'>🚨 触发中告警 [1]</font>",
		"<font color='green'>✅ 已恢复告警 [1]</font>",
	} {
		idx := strings.Index(out, title)
		if idx < 0 {
			t.Fatalf("missing %q in output:\n%s", title, out)
		}
		after := out[idx+len(title):]
		if !strings.HasPrefix(after, "\n\n\n") {
			t.Fatalf("expected two blank lines after summary %q, got:\n%s", title, out)
		}
	}
}
