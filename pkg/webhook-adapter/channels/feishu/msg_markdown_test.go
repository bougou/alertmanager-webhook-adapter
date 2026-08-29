package feishu

import (
	"testing"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func TestHeaderColor(t *testing.T) {
	cases := []struct {
		severity string
		want     string
	}{
		{severity: models.SeverityCritical, want: "red"},
		{severity: models.SeverityWarning, want: "orange"},
		{severity: models.SeverityInfo, want: "blue"},
		{severity: models.SeverityOK, want: "green"},
		{severity: "", want: "turquoise"},
	}
	for _, tc := range cases {
		if got := headerColor(tc.severity); got != tc.want {
			t.Errorf("headerColor(%q) = %q, want %q", tc.severity, got, tc.want)
		}
	}
}

func TestNewMsgMarkdownFromPayloadSetsHeaderTemplate(t *testing.T) {
	msg := NewMsgMarkdownFromPayload(&models.Payload{
		Title:    "[ops] WARNING • CPUHigh • Firing:1",
		Markdown: "**Summary**",
		Status:   models.StatusFiring,
		Severity: models.SeverityWarning,
	})
	if msg.Card == nil || msg.Card.Header == nil {
		t.Fatal("expected card header")
	}
	if msg.Card.Header.Template != "orange" {
		t.Fatalf("header template = %q, want orange", msg.Card.Header.Template)
	}
}

func TestNewMsgMarkdownFromPayloadResolvedUsesGreen(t *testing.T) {
	msg := NewMsgMarkdownFromPayload(&models.Payload{
		Title:    "[ops] CRITICAL • CPUHigh • Resolved:1",
		Markdown: "**Summary**",
		Status:   models.StatusResolved,
		Severity: models.SeverityCritical,
	})
	if msg.Card.Header.Template != "green" {
		t.Fatalf("resolved header template = %q, want green", msg.Card.Header.Template)
	}
}

func TestNewCardMarkdownSetsHeaderTemplate(t *testing.T) {
	msg := NewMsgMarkdown("[ops] WARNING • CPUHigh • Firing:1", "**Summary**")
	if msg.Card == nil || msg.Card.Header == nil {
		t.Fatal("expected card header")
	}
	if msg.Card.Header.Template != "turquoise" {
		t.Fatalf("header template without payload status = %q, want turquoise", msg.Card.Header.Template)
	}
}
