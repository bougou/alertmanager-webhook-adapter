package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// WexinGroupBot can send messages to weixin group
// ref: https://work.weixin.qq.com/api/doc/90000/90136/91770
type WeixinGroupBot struct {
	addr   string
	key    string
	client *http.Client
}

func NewWexinGroupBot(key string) *WeixinGroupBot {
	return &WeixinGroupBot{
		addr:   "https://qyapi.weixin.qq.com",
		key:    key,
		client: &http.Client{},
	}
}

func (bot *WeixinGroupBot) Addr() string {
	return fmt.Sprintf("%s/cgi-bin/webhook/send?key=%s", bot.addr, bot.key)
}

func (bot *WeixinGroupBot) AddrForUpload() string {
	return fmt.Sprintf("%s/cgi-bin/webhook/upload_media?key=%s&type=file", bot.addr, bot.key)
}

func (bot *WeixinGroupBot) Send(msg *Msg) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", bot.Addr(), bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("failed to construct request, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := bot.client.Do(req)
	if err != nil {
		return fmt.Errorf("send msg error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("send msg response error, status: %d", res.StatusCode)
	}

	type WeixinGroupBotResponse struct {
		Errcode int    `json:"errcode,omitempty"`
		Errmsg  string `json:"errmsg,omitempty"`
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read res body failed, err: %s", err)
	}

	botRes := &WeixinGroupBotResponse{}
	if err := json.Unmarshal(resBody, botRes); err != nil {
		return fmt.Errorf("marshal to WeixinGroupBotResponse failed, err: %s", err)
	}

	if botRes.Errcode != 0 {
		return fmt.Errorf("found err in response, request url: (%s), response: %s", bot.Addr(), string(resBody))
	}

	return nil
}

func (bot *WeixinGroupBot) UploadFile(filename string, fileReader io.Reader) (meidaID string, err error) {

	// Todo
	// 要求文件大小在5B~20M之间

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	fileWriter, err := w.CreateFormFile("media", filename)
	if err != nil {
		return "", fmt.Errorf("create file writer error, got %v", err)
	}
	io.Copy(fileWriter, fileReader)

	req, err := http.NewRequest("POST", bot.AddrForUpload(), body)
	if err != nil {
		return "", fmt.Errorf("failed to construct request, got %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := bot.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("upload file failed, status: %d", res.StatusCode)
	}

	type ResponseBody struct {
		ErrCode   int       `json:"errcode"`
		ErrMsg    string    `json:"errmsg"`
		Type      string    `json:"type"`       // 媒体文件类型，分别有图片（image）、语音（voice）、视频（video），普通文件(file)
		MediaID   string    `json:"media_id"`   // 媒体文件上传后获取的唯一标识，3天内有效
		CreatedAt time.Time `json:"created_at"` // 媒体文件上传时间戳
	}

	r := &ResponseBody{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return "", fmt.Errorf("can decode res body, err: %v", err)
	}

	return r.MediaID, nil
}
