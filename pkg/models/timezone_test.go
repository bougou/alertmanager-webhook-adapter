package models

import (
	"os"
	"testing"
	"time"
)

// testCST matches the timezone used by golden files. Template rendering calls
// time.Time.Format, which uses the time's location; Alert.UnmarshalJSON maps
// timestamps into time.Local. Pinning Local keeps tests independent of the host TZ.
var testCST = time.FixedZone("CST", 8*3600)

func TestMain(m *testing.M) {
	time.Local = testCST
	os.Exit(m.Run())
}
