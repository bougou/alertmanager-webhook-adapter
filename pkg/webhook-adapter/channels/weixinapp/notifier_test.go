package weixinapp

import (
	"os"
	"strconv"
	"testing"
)

func Test_Notifier_Send(t *testing.T) {
	coreID := os.Getenv("WEIXIN_APP_CORP_ID")

	aID := os.Getenv("WEIXIN_APP_AGENT_ID")
	agentID, err := strconv.Atoi(aID)
	if err != nil {
		t.Fatal(err)
	}

	agentSecret := os.Getenv("WEIXIN_APP_SECRET")
	toUser := ""
	toParty := "2"
	toTag := ""
	tt := NewNotifer(coreID, agentID, agentSecret, toUser, toParty, toTag)

	msg := NewMsgMarkdown("# Hello World")
	if err := tt.Send(msg); err != nil {
		t.Fatalf("err: %v", err)
	}
}
