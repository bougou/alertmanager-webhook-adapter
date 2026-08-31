package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func loadWebhookJSON(t *testing.T, rel string) *AlertmanagerWebhookMessage {
	t.Helper()
	b, err := os.ReadFile(rel)
	if err != nil {
		t.Fatal(err)
	}
	var msg AlertmanagerWebhookMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		t.Fatal(err)
	}
	msg.Signature = "Bougou"
	msg.MessageAt = time.Date(2026, 8, 29, 17, 31, 23, 0, testCST)
	for i := range msg.Alerts {
		msg.Alerts[i].StartsAt = msg.Alerts[i].StartsAt.In(testCST)
		msg.Alerts[i].EndsAt = msg.Alerts[i].EndsAt.In(testCST)
	}
	return &msg
}

func renderFeishuTmpl(t *testing.T, src string, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	tmpl, err := template.New("prom").
		Funcs(defaultFuncs).
		Option("missingkey=zero").
		Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, name, msg); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func renderFeishuZH(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplFeishuZH, msg, name)
}

func renderFeishuEN(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplFeishu, msg, name)
}

func TestFeishuZHTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderFeishuZH(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderFeishuZH(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "feishu_zh.golden")
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
		t.Fatalf("feishu zh render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	md := renderFeishuZH(t, msg, "prom.markdown")
	if !strings.HasSuffix(md, "\n") || strings.HasSuffix(md, "\n\n") {
		t.Fatalf("markdown should end with exactly one newline, got trailing %q", md[len(md)-8:])
	}
}

func TestFeishuZHMarkdownSummaryIsCompactList(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderFeishuZH(t, msg, "prom.markdown")
	summary := out[strings.Index(out, "**摘要**"):strings.Index(out, "**详情**")]
	if strings.Contains(summary, "</text_tag>\n\n- ") {
		t.Fatalf("summary list items should be consecutive, got:\n%s", summary)
	}
	if !strings.Contains(summary, "</text_tag>\n- ") {
		t.Fatalf("summary list items should be separated by a single newline, got:\n%s", summary)
	}
}

func TestFeishuZHDetailUsesTwoBlankLines(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	out := renderFeishuZH(t, msg, "prom.markdown")
	detail := out[strings.Index(out, "**详情**"):]
	title := "<font color='red'>🚨 触发中告警 [3]</font>"
	idx := strings.Index(detail, title)
	if idx < 0 {
		t.Fatalf("missing firing title in detail, got:\n%s", detail)
	}
	afterTitle := detail[idx+len(title):]
	if !strings.HasPrefix(afterTitle, "\n\n\n- **告警实例**:") {
		t.Fatalf("detail title should be followed by two blank lines, got:\n%s", afterTitle[:min(80, len(afterTitle))])
	}
	if strings.HasPrefix(afterTitle, "\n\n\n\n") {
		t.Fatalf("detail title has more than two blank lines, got:\n%s", afterTitle[:min(80, len(afterTitle))])
	}

	first := strings.Index(detail, "- **告警实例**: <text_tag color='purple'>10.30.1.160</text_tag>")
	second := strings.Index(detail, "- **告警实例**: <text_tag color='purple'>10.30.1.162</text_tag>")
	if first < 0 || second < 0 {
		t.Fatalf("missing firing detail items, got:\n%s", detail)
	}
	gap := detail[first:second]
	if !strings.HasSuffix(gap, "\n\n\n") {
		t.Fatalf("detail items should be separated by two blank lines, gap ending %q", gap[len(gap)-12:])
	}
	if strings.HasSuffix(gap, "\n\n\n\n") {
		t.Fatalf("detail items have more than two blank lines, gap ending %q", gap[len(gap)-16:])
	}
}

func TestFeishuTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderFeishuEN(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderFeishuEN(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "feishu.golden")
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
		t.Fatalf("feishu en render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestFeishuENLayoutMatchesZH(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	zhTitle := renderFeishuZH(t, msg, "prom.title")
	enTitle := renderFeishuEN(t, msg, "prom.title")
	zhMD := renderFeishuZH(t, msg, "prom.markdown")
	enMD := renderFeishuEN(t, msg, "prom.markdown")
	if feishuZHToEN(zhTitle) != enTitle {
		t.Fatalf("title layout mismatch\nzh-mapped:\n%s\nen:\n%s", feishuZHToEN(zhTitle), enTitle)
	}
	if feishuZHToEN(zhMD) != enMD {
		t.Fatalf("markdown layout mismatch\nzh-mapped:\n%s\nen:\n%s", feishuZHToEN(zhMD), enMD)
	}
}

func feishuZHToEN(s string) string {
	s = regexp.MustCompile(`触发中告警 \[(\d+)\]`).ReplaceAllString(s, "Firing [$1] alerts")
	s = regexp.MustCompile(`已恢复告警 \[(\d+)\]`).ReplaceAllString(s, "Resolved [$1] alerts")
	replacer := strings.NewReplacer(
		"**摘要**", "**Summary**",
		"**详情**", "**Detail**",
		"**告警名称**", "**Alert Name**",
		"**告警级别**", "**Alert Level**",
		"**告警状态**", "**Alert Status**",
		"**告警描述**", "**Description**",
		"**告警实例**", "**Alert Instance**",
		"**可用区**", "**Zone**",
		"**地域**", "**Region**",
		"**产品**", "**Product**",
		"**组件**", "**Component**",
		"**开始**", "**Start At**",
		"**结束**", "**End At**",
		"告警中:", "Firing:",
		"已恢复:", "Resolved:",
		"未结束", "Not End",
	)
	return replacer.Replace(s)
}
