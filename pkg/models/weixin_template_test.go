package models

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func renderWeixin(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplWeixin, msg, name)
}

func TestWeixinTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderWeixin(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderWeixin(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "weixin.golden")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(wantPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(wantPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("weixin render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	md := renderWeixin(t, msg, "prom.markdown")
	if !strings.HasSuffix(md, "\n") || strings.HasSuffix(md, "\n\n") {
		t.Fatalf("markdown should end with exactly one newline, got trailing %q", md[len(md)-8:])
	}
}

func TestWeixinTitleUsesHighestSeverity(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{Status: "firing", Labels: KV{"severity": "warning", "alertname": "CPUHigh"}},
			{Status: "firing", Labels: KV{"severity": "critical", "alertname": "CPUHigh"}},
		},
	}
	out := renderWeixin(t, msg, "prom.title")
	if !strings.Contains(out, "CRITICAL") {
		t.Fatalf("title should use highest severity CRITICAL, got: %s", out)
	}
	if strings.Contains(out, "WARNING") {
		t.Fatalf("title should not keep first-alert WARNING, got: %s", out)
	}
}

func TestWeixinENLayoutMatchesZH(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	if weixinZHToEN(renderWeixinZH(t, msg, "prom.title")) != renderWeixin(t, msg, "prom.title") {
		t.Fatalf("title layout mismatch\nzh-mapped:\n%s\nen:\n%s",
			weixinZHToEN(renderWeixinZH(t, msg, "prom.title")), renderWeixin(t, msg, "prom.title"))
	}
	if weixinZHToEN(renderWeixinZH(t, msg, "prom.markdown")) != renderWeixin(t, msg, "prom.markdown") {
		t.Fatalf("markdown layout mismatch\nzh-mapped:\n%s\nen:\n%s",
			weixinZHToEN(renderWeixinZH(t, msg, "prom.markdown")), renderWeixin(t, msg, "prom.markdown"))
	}
}

func weixinZHToEN(s string) string {
	s = regexp.MustCompile(`触发中告警 \[(\d+)\]`).ReplaceAllString(s, "Firing Alerts [$1]")
	s = regexp.MustCompile(`已恢复告警 \[(\d+)\]`).ReplaceAllString(s, "Resolved Alerts [$1]")
	return strings.NewReplacer(
		"告警名称", "Alert Name",
		"告警级别", "Alert Level",
		"告警状态", "Alert Status",
		"告警描述", "Description",
		"未结束", "Not End",
		"告警中:", "Firing:",
		"已恢复:", "Resolved:",
		"#### 摘要", "#### Summary",
		"#### 详情", "#### Detail",
		"实例", "Instance",
		"可用区", "Zone",
		"地域", "Region",
		"产品", "Product",
		"组件", "Component",
		"开始", "Start At",
		"结束", "End At",
	).Replace(s)
}
