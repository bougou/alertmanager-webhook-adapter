package models

import (
	"os"
	"testing"
	"time"
)

// testUTC8 pins timestamps to UTC+8 so golden files stay stable. Template
// rendering uses numeric offsets (e.g. +08); Alert.UnmarshalJSON maps times
// into time.Local. Pinning Local keeps tests independent of the host TZ.
var testCST = time.FixedZone("CST", 8*3600)

func TestMain(m *testing.M) {
	time.Local = testCST
	os.Exit(m.Run())
}
