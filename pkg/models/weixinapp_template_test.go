package models

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func renderWeixinapp(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplWeixinapp, msg, name)
}

func TestWeixinappTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderWeixinapp(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderWeixinapp(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "weixinapp.golden")
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
		t.Fatalf("weixinapp render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestWeixinappENLayoutMatchesZH(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	if weixinappZHToEN(renderWeixinappZH(t, msg, "prom.title")) != renderWeixinapp(t, msg, "prom.title") {
		t.Fatalf("title layout mismatch\nzh-mapped:\n%s\nen:\n%s",
			weixinappZHToEN(renderWeixinappZH(t, msg, "prom.title")), renderWeixinapp(t, msg, "prom.title"))
	}
	if weixinappZHToEN(renderWeixinappZH(t, msg, "prom.markdown")) != renderWeixinapp(t, msg, "prom.markdown") {
		t.Fatalf("markdown layout mismatch\nzh-mapped:\n%s\nen:\n%s",
			weixinappZHToEN(renderWeixinappZH(t, msg, "prom.markdown")), renderWeixinapp(t, msg, "prom.markdown"))
	}
}

func weixinappZHToEN(s string) string {
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
		"### 告警摘要", "### Alert Summary",
		"### 告警详请", "### Alert Detail",
		"告警实例", "Alert Instance",
		"可用区", "Zone",
		"地域", "Region",
		"产品", "Product",
		"组件", "Component",
		"开始", "Start At",
		"结束", "End At",
	).Replace(s)
}
