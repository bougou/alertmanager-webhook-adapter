package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func renderWeixinZH(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplWeixinZH, msg, name)
}

func TestWeixinZHTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderWeixinZH(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderWeixinZH(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "weixin_zh.golden")
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
		t.Fatalf("weixin zh render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	md := renderWeixinZH(t, msg, "prom.markdown")
	if !strings.HasSuffix(md, "\n") || strings.HasSuffix(md, "\n\n") {
		t.Fatalf("markdown should end with exactly one newline, got trailing %q", md[len(md)-8:])
	}
}

func TestWeixinZHTitleUsesHighestSeverity(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{Status: "firing", Labels: KV{"severity": "warning", "alertname": "CPUHigh"}},
			{Status: "firing", Labels: KV{"severity": "critical", "alertname": "CPUHigh"}},
		},
	}
	out := renderWeixinZH(t, msg, "prom.title")
	if !strings.Contains(out, "CRITICAL") {
		t.Fatalf("title should use highest severity CRITICAL, got: %s", out)
	}
	if strings.Contains(out, "WARNING") {
		t.Fatalf("title should not keep first-alert WARNING, got: %s", out)
	}
}

func TestWeixinZHMarkdownSummaryIsPlainInstanceList(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderWeixinZH(t, msg, "prom.markdown")
	summary := out[strings.Index(out, "#### 摘要"):strings.Index(out, "#### 详情")]
	if strings.Contains(summary, "- <font") {
		t.Fatalf("summary items should not use markdown list markers, got:\n%s", summary)
	}
	if !strings.Contains(summary, "• <font") {
		t.Fatalf("summary items should use bullet prefix, got:\n%s", summary)
	}
	if strings.Contains(summary, "10.30.1.160</font>\n\n• <font") {
		t.Fatalf("summary list items should be consecutive, got:\n%s", summary)
	}
	if !strings.Contains(summary, "10.30.1.160</font>\n• <font") {
		t.Fatalf("summary list items should be separated by a single newline, got:\n%s", summary)
	}
}

func TestWeixinZHDetailUsesWeixinMarkdown(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderWeixinZH(t, msg, "prom.markdown")
	if strings.Contains(out, "<text_tag") {
		t.Fatalf("weixin markdown must not use feishu text_tag")
	}
	if strings.Contains(out, "color=\"#") || strings.Contains(out, "color='#") {
		t.Fatalf("weixin markdown only supports info/comment/warning colors")
	}
	if strings.Contains(out, "\n---\n") {
		t.Fatalf("weixin classic markdown does not support horizontal rules")
	}
	if strings.Contains(out, "详请") {
		t.Fatalf("should use 详情 not 详请")
	}
	if strings.Contains(out, "Not End") {
		t.Fatalf("zh template should use 未结束")
	}

	detail := out[strings.Index(out, "#### 详情"):]
	title := `<font color="warning">触发中告警 [3]</font>`
	idx := strings.Index(detail, title)
	if idx < 0 {
		t.Fatalf("missing firing title in detail, got:\n%s", detail)
	}
	afterTitle := detail[idx+len(title):]
	if !strings.HasPrefix(afterTitle, "\n\n\n> <font color=\"comment\">实例</font>:") {
		t.Fatalf("detail title should be followed by two blank lines then a quote, got:\n%s", afterTitle[:min(80, len(afterTitle))])
	}

	first := strings.Index(detail, "> <font color=\"comment\">实例</font>: `10.30.1.160`")
	second := strings.Index(detail, "> <font color=\"comment\">实例</font>: `10.30.1.162`")
	if first < 0 || second < 0 {
		t.Fatalf("missing firing detail items, got:\n%s", detail)
	}
	gap := detail[first:second]
	if strings.Contains(gap, "\n>\n>") {
		t.Fatalf("weixin quotes should be consecutive without empty > lines, gap:\n%s", gap)
	}
	if !strings.HasSuffix(gap, "\n\n\n") {
		t.Fatalf("detail items should be separated by two blank lines, gap ending %q", gap[len(gap)-12:])
	}
}
