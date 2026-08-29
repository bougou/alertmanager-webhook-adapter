package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

func renderWeixinappZH(t *testing.T, msg *AlertmanagerWebhookMessage, name string) string {
	t.Helper()
	return renderFeishuTmpl(t, templates.DefaultTmplWeixinappZH, msg, name)
}

func TestWeixinappZHTemplateGolden(t *testing.T) {
	msg := loadWebhookJSON(t, filepath.Join("..", "..", "tests", "alert.json"))
	got := "TITLE:\n" + renderWeixinappZH(t, msg, "prom.title") + "\nMARKDOWN:\n" + renderWeixinappZH(t, msg, "prom.markdown")
	wantPath := filepath.Join("testdata", "weixinapp_zh.golden")
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
		t.Fatalf("weixinapp zh render changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
