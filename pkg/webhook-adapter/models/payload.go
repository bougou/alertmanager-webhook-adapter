package models

import "strings"

const (
	StatusFiring   = "firing"
	StatusResolved = "resolved"

	SeverityCritical = "critical"
	SeverityWarning  = "warning"
	SeverityInfo     = "info"
	SeverityOK       = "ok"
)

// The purpse of Payload is to hide the complexity of constructing channel-specific Msg.
// Because each specific channel provides suitable Payload2MsgFn(s) convertion functions for its supported msgType(s).
type Payload struct {
	Raw      string   `json:"raw"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Markdown string   `json:"markdown"` // Don't put Title content in Markdown
	Files    []string `json:"files"`
	Images   []Image  `json:"images"`
	Links    []Link   `json:"links"`
	Buttons  []Button `json:"buttons"`
	At       At       `json:"at"`

	// Status and Severity come from the Alertmanager webhook, not from rendered title text.
	Status   string `json:"status"`   // firing | resolved
	Severity string `json:"severity"` // critical | warning | info | ok
}

type PayloadGenerator interface {
	ToPayload() *Payload
}

type Image struct {
	Bytes  []byte `json:"bytes"`
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

type Link struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	PicURL string `json:"picURL"`
}

type Button struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtAll     bool     `json:"atAll"`
}

// EffectiveSeverity is the severity channels should use for coloring and emphasis.
// A fully resolved group is treated as ok, even if original labels were higher.
func (p *Payload) EffectiveSeverity() string {
	if p == nil {
		return ""
	}
	if strings.EqualFold(p.Status, StatusResolved) {
		return SeverityOK
	}
	switch s := strings.ToLower(p.Severity); s {
	case SeverityCritical, SeverityWarning, SeverityInfo, SeverityOK:
		return s
	default:
		return ""
	}
}
