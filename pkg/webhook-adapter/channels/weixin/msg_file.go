package weixin

import (
	"fmt"
	"io"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeFile] = NewMsgFileFromPayload
}

type File struct {
	MediaID string `json:"media_id"`
}

func NewMsgFile(file *File) *Msg {
	return &Msg{
		MsgType: MsgTypeFile,
		File:    file,
	}
}

func NewMsgFileFromMediaID(mediaID string) *Msg {
	return &Msg{
		MsgType: MsgTypeFile,
		File: &File{
			MediaID: mediaID,
		},
	}
}

func NewMsgFileFromPayload(payload *models.Payload) *Msg {
	// Todo, first upload file to get mediaID
	mediaID := ""
	return NewMsgFileFromMediaID(mediaID)
}

func (b *WeixinGroupBot) SendFile(filename string, fileReader io.Reader) error {
	mediaID, err := b.UploadFile(filename, fileReader)
	if err != nil {
		return fmt.Errorf("send file error, err: %v", err)
	}

	msg := NewMsgFileFromMediaID(mediaID)
	return b.Send(msg)
}
