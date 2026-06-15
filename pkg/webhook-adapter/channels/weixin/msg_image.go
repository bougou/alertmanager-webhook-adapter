package weixin

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/webhook-adapter/models"
)

func init() {
	Payload2MsgFnMap[MsgTypeImage] = NewMsgImageFromPayload
}

type Image struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

func NewImageFromBytes(imgByte []byte) *Image {
	imgMD5 := GetMD5Hash(imgByte)
	imgBase64 := base64.StdEncoding.EncodeToString(imgByte)

	return &Image{
		Base64: imgBase64,
		MD5:    string(imgMD5),
	}
}

func GetMD5Hash(data []byte) []byte {
	md5sum := md5.Sum(data)
	return []byte(hex.EncodeToString(md5sum[:]))
}

func NewMsgImageFromBytes(imgByte []byte) *Msg {
	return &Msg{
		MsgType: MsgTypeImage,
		Image:   NewImageFromBytes(imgByte),
	}
}

func NewMsgImage(image *Image) *Msg {
	return &Msg{
		MsgType: MsgTypeImage,
		Image:   image,
	}
}

func NewMsgImageFromPayload(payload *models.Payload) *Msg {
	imgByte := []byte{}
	if len(payload.Images) > 0 {
		imgByte = payload.Images[0].Bytes
	}
	return NewMsgImageFromBytes(imgByte)
}
