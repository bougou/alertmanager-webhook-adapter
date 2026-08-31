package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func renderDingtalkZH(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplDingTalkZH, msg, name)
}

func TestDingtalkZHTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderDingtalkZH(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderDingtalkZH(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "dingtalk_zh.golden")
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
		t.Fatalf("dingtalk zh render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	md := renderDingtalkZH(t, msg, "prom.markdown")
	if !strings.HasSuffix(md, ">\n") || strings.HasSuffix(md, ">\n\n") {
		t.Fatalf("markdown should end with a quote line and exactly one newline, got trailing %q", md[len(md)-8:])
	}
}

func TestDingtalkZHTitleUsesHighestSeverity(t *testing.T) {
	msg := &AlertmanagerWebhookMessage{
		Signature:   "ops",
		GroupLabels: KV{"alertname": "CPUHigh"},
		Alerts: Alerts{
			{Status: "firing", Labels: KV{"severity": "warning", "alertname": "CPUHigh"}},
			{Status: "firing", Labels: KV{"severity": "critical", "alertname": "CPUHigh"}},
		},
	}
	out := renderDingtalkZH(t, msg, "prom.title")
	if !strings.Contains(out, "CRITICAL") {
		t.Fatalf("title should use highest severity CRITICAL, got: %s", out)
	}
	if strings.Contains(out, "WARNING") {
		t.Fatalf("title should not keep first-alert WARNING, got: %s", out)
	}
}

func TestDingtalkZHMarkdownSummaryIsCompactList(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderDingtalkZH(t, msg, "prom.markdown")
	summary := out[strings.Index(out, "**摘要**"):strings.Index(out, "**详情**")]
	if strings.Contains(summary, "10.30.1.160\n\n- ") {
		t.Fatalf("summary list items should be consecutive, got:\n%s", summary)
	}
	if !strings.Contains(summary, "10.30.1.160\n- ") {
		t.Fatalf("summary list items should be separated by a single newline, got:\n%s", summary)
	}
	if strings.Contains(summary, "> <font") {
		t.Fatalf("summary should use unordered list, not blockquote, got:\n%s", summary)
	}
}

func TestDingtalkZHDetailSeparatesAlerts(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderDingtalkZH(t, msg, "prom.markdown")
	detail := out[strings.Index(out, "**详情**"):]
	title := `<font color="#FF4D4F">**🚨 触发中告警 [3]**</font>`
	idx := strings.Index(detail, title)
	if idx < 0 {
		t.Fatalf("missing firing title in detail, got:\n%s", detail)
	}
	afterTitle := detail[idx+len(title):]
	if !strings.Contains(afterTitle, "\n> **告警实例**:") {
		t.Fatalf("detail title should be followed by a quoted alert, got:\n%s", afterTitle[:min(80, len(afterTitle))])
	}

	first := strings.Index(detail, `> **告警实例**: <font color="#722ED1">10.30.1.160</font>`)
	second := strings.Index(detail, `> **告警实例**: <font color="#722ED1">10.30.1.162</font>`)
	if first < 0 || second < 0 {
		t.Fatalf("missing firing detail items, got:\n%s", detail)
	}
	gap := detail[first:second]
	if !strings.Contains(gap, "\n\n---\n\n") {
		t.Fatalf("quoted alerts should be separated by a blank line, ---, and a blank line, gap ending %q", gap[len(gap)-24:])
	}
	if strings.Contains(gap, "\n- **告警实例**") {
		t.Fatalf("detail alerts should use blockquote, not unordered list, gap:\n%s", gap)
	}
	if strings.Contains(gap, "**告警实例**:  \n<font") {
		t.Fatalf("key and value should stay on the same line, gap:\n%s", gap)
	}
	if !strings.Contains(gap, "\n>\n> **告警名称**:") {
		t.Fatalf("fields of one alert should stay in one quote block with empty > line breaks, gap:\n%s", gap)
	}
}

func TestDingtalkZHSectionHeadingsHaveBlankLines(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderDingtalkZH(t, msg, "prom.markdown")
	for _, want := range []string{
		"**摘要**\n\n",
		"**详情**\n\n",
		"<font color=\"#FF4D4F\">**🚨 触发中告警 [3]**</font>\n\n",
		"<font color=\"#52C41A\">**✅ 已恢复告警 [2]**</font>\n\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing blank lines around section headings, want %q, got:\n%s", want, out)
		}
	}
}

func TestDingtalkZHUsesOfficialLineBreaks(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderDingtalkZH(t, msg, "prom.markdown")
	if strings.Contains(out, "详请") {
		t.Fatalf("should use 详情 not 详请")
	}
	if strings.Contains(out, "Not End") {
		t.Fatalf("zh template should use 未结束")
	}
}
