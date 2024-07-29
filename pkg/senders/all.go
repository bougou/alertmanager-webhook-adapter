package senders

import (
	"fmt"

	prommodels "github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

const (
	TimeFormat string = "2006-01-02 15:04:05"
)

func ErrNotFoundPayload2MsgFn(channelType string, msgType string) error {
	return fmt.Errorf("not found payload2MsgFn for channel_type/msg_type (%s/%s)", channelType, msgType)
}

func ErrNotFoundConverter(channelType string, msgType string) error {
	return fmt.Errorf("not found msg converter for channel_type/msg_type (%s/%s)", channelType, msgType)
}

// ChannelSenderCreator parses the required arguments for creating sender
// like (token, msg_type, ...) from the http request, and returns the created sender
// and the msg converter for corresponding the msg_type.
type ChannelSenderCreator func(request *restful.Request) (models.Sender, MsgConverter, error)

// ChannelsSenderCreatorMap holds a registry for SenderCreator for each channel type.
//
// the key of the map is the ChannelType
var ChannelsSenderCreatorMap map[string]ChannelSenderCreator

func RegisterChannelsSenderCreator(channel string, creator ChannelSenderCreator) {
	if ChannelsSenderCreatorMap == nil {
		ChannelsSenderCreatorMap = make(map[string]ChannelSenderCreator)
	}
	ChannelsSenderCreatorMap[channel] = creator
}

type MsgConverter interface {
	Convert(raw []byte, promMsg *prommodels.AlertmanagerWebhookMessage) (interface{}, error)
}

// map[ChannelType]map[MsgType]MsgConverter
var ChannelsMsgConverterMap map[string]map[string]MsgConverter

func RegisterChannelsMsgConverter(channelType string, msgType string, msgConverter MsgConverter) {
	if ChannelsMsgConverterMap == nil {
		ChannelsMsgConverterMap = make(map[string]map[string]MsgConverter)
	}

	m, exist := ChannelsMsgConverterMap[channelType]
	if !exist {
		m = make(map[string]MsgConverter)
		ChannelsMsgConverterMap[channelType] = m
	}

	m[msgType] = msgConverter
}

func getMsgConverter(channelType string, msgType string) (converter MsgConverter, exists bool) {
	m, exists := ChannelsMsgConverterMap[channelType]
	if !exists {
		return nil, false
	}

	converter, exists = m[msgType]
	if !exists {
		return nil, false
	}

	return converter, true
}
