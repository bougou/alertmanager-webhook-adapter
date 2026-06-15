package weixinapp

import (
	"fmt"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

// 企业微信 - 应用

const ChannelTypeWeixin = "weixinapp"

const (
	MsgTypeFile string = "file"

	MsgTypeImage string = "image"

	MsgTypeInteractiveCard string = "interactive_card"

	MsgTypeMarkdown  string = "markdown"
	maxMarkdownBytes int    = 2048
	truncatedMark    string = "\n... more is truncated due to limit"
	// see: https://work.weixin.qq.com/api/doc/90000/90135/90236#markdown%E6%B6%88%E6%81%AF

	MsgTypeMiniProgramNotice string = "miniprogram_notice"

	MsgTypeMPNews string = "mpnews"

	MsgTypeNews         string = "news"
	maxArticlesNumber   int    = 8
	maxTitleBytes       int    = 128
	maxDescriptionBytes int    = 512

	MsgTypeText string = "text"

	MsgTypeTextCard string = "textcard"

	MsgTypeVideo string = "video"

	MsgTypeVoice string = "voice"
)

type Payload2MsgFn func(payload *models.Payload) *Msg

var Payload2MsgFnMap = make(map[string]Payload2MsgFn)

type Msg struct {
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`

	// touser、toparty、totag不能同时为空
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`

	Text                *Text                `json:"text,omitempty"`
	Image               *Image               `json:"image,omitempty"`
	Voice               *Voice               `json:"voice,omitempty"`
	File                *File                `json:"file,omitempty"`
	TextCard            *TextCard            `json:"textcard,omitempty"`
	News                *News                `json:"news,omitempty"`
	MPNews              *MPNews              `json:"mpnews,omitempty"`
	Markdown            *Markdown            `json:"markdown,omitempty"`
	MiniprogramNotice   *MiniprogramNotice   `json:"miniprogram_notice,omitempty"`
	InteractiveTaskcard *InteractiveTaskcard `json:"interactive_taskcard,omitempty"`

	Safe                   int `json:"safe,omitempty"`
	EnableIDTrans          int `json:"enable_id_trans,omitempty"`
	DuplicateCheckInterval int `json:"duplicate_check_interval,omitempty"`
}

func validateMsg(msgType string, msg *Msg) error {
	if msg.MsgType != msgType {
		return fmt.Errorf("the msg does not match with specified msgType")
	}

	switch msgType {
	case MsgTypeFile:
	case MsgTypeImage:
	case MsgTypeMarkdown:
	case MsgTypeNews:
	case MsgTypeText:
	default:
		return fmt.Errorf("unsupported msgtype of (%s)", msgType)
	}

	return nil
}
