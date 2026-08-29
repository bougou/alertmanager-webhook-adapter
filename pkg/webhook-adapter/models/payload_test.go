package models

import "testing"

func TestPayloadEffectiveSeverity(t *testing.T) {
	cases := []struct {
		status   string
		severity string
		want     string
	}{
		{status: StatusFiring, severity: SeverityCritical, want: SeverityCritical},
		{status: StatusFiring, severity: SeverityWarning, want: SeverityWarning},
		{status: StatusResolved, severity: SeverityWarning, want: SeverityOK},
		{status: StatusResolved, severity: SeverityCritical, want: SeverityOK},
		{status: StatusFiring, severity: SeverityInfo, want: SeverityInfo},
		{status: StatusResolved, severity: SeverityOK, want: SeverityOK},
		{status: StatusFiring, severity: "", want: ""},
		{status: "FIRING", severity: "CRITICAL", want: SeverityCritical},
	}
	for _, tc := range cases {
		p := &Payload{Status: tc.status, Severity: tc.severity}
		if got := p.EffectiveSeverity(); got != tc.want {
			t.Errorf("EffectiveSeverity(%q, %q) = %q, want %q", tc.status, tc.severity, got, tc.want)
		}
	}
}
