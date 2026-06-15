package senders

import (
	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
	restful "github.com/emicklei/go-restful/v3"
)

type ChannelSenderCreator func(request *restful.Request) (models.Sender, error)

var ChannelsSenderCreatorMap = map[string]ChannelSenderCreator{}

func RegisterChannelsSenderCreator(channel string, creator ChannelSenderCreator) {
	ChannelsSenderCreatorMap[channel] = creator
}
