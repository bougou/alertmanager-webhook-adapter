package feishu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// FeishuBot can send messages to feishu group
// ref: https://www.feishu.cn/hc/zh-CN/articles/360024984973
type FeishuGroupBot struct {
	addr   string
	token  string
	client *http.Client
}

func NewFeishuGroupBot(token string) *FeishuGroupBot {
	return &FeishuGroupBot{
		addr:   "https://open.feishu.cn",
		token:  token,
		client: &http.Client{},
	}
}

func (bot *FeishuGroupBot) Addr() string {
	// Note, use v2
	return fmt.Sprintf("%s/open-apis/bot/v2/hook/%s", bot.addr, bot.token)
}

func (bot *FeishuGroupBot) AddrForUploadImage() string {
	return fmt.Sprintf("%s/open-apis/image/v4/put", bot.addr)
}

func (bot *FeishuGroupBot) AddrForFetchImage() string {
	return fmt.Sprintf("%s/open-apis/image/v4/get", bot.addr)
}

func (bot *FeishuGroupBot) AddrForUploadFile() string {
	return fmt.Sprintf("%s/open-apis/im/v1/files", bot.addr)
}

func (bot *FeishuGroupBot) AddrForDownloadFile(fileKey string) string {
	return fmt.Sprintf("%s/open-apis/im/v1/files/%s", bot.addr, fileKey)
}

func (bot *FeishuGroupBot) Send(msg *Msg) error {
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

	return nil
}

func GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (bot *FeishuGroupBot) UploadImage(filename string, fileReader io.Reader) (imageKey string, err error) {

	postBody := new(bytes.Buffer)
	w := multipart.NewWriter(postBody)
	fileWriter, err := w.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("create form file err")
	}
	io.Copy(fileWriter, fileReader)

	fieldWriter, err := w.CreateFormField("image_type")
	if err != nil {
		return "", fmt.Errorf("create field error")
	}

	fieldWriter.Write([]byte("message")) // message 表示消息图片， avatar 表示头像

	req, err := http.NewRequest("POST", bot.AddrForUploadImage(), postBody)
	if err != nil {
		return "", fmt.Errorf("failed to construct request")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.token))

	res, err := bot.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("upload file failed, status: %d", res.StatusCode)
	}

	type ResponseBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ImageKey string `json:"image_key"` // 媒体文件上传后获取的唯一标识，3天内有效
		} `json:"data"`
	}

	r := &ResponseBody{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return "", fmt.Errorf("can decode res body, err: %v", err)
	}

	return r.Data.ImageKey, nil
}

func (bot *FeishuGroupBot) FetchImage(imageKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", bot.AddrForFetchImage(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request")
	}

	q := req.URL.Query()
	q.Add("image_key", imageKey)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.token))

	res, err := bot.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("fetch image failed, status: %d", res.StatusCode)
	}

	result := []byte{}
	if _, err := res.Body.Read(result); err != nil {
		return nil, fmt.Errorf("read image body failed")
	}

	return result, nil
}

func (bot *FeishuGroupBot) UploadFile(filename string, filetype string, fileReader io.Reader) (fileKey string, err error) {

	postBody := new(bytes.Buffer)
	w := multipart.NewWriter(postBody)
	var field io.Writer

	// filetype, 文件类型,可选的类型有mp4,pdf,doc
	field, err = w.CreateFormField("file_type")
	if err != nil {
		return "", fmt.Errorf("create field error, got %v", err)
	}
	field.Write([]byte(filetype))

	// 带后缀的文件名
	field, err = w.CreateFormField("file_name")
	if err != nil {
		return "", fmt.Errorf("create field error, got %v", err)
	}
	field.Write([]byte("filename"))

	// 文件的时长(视频，音频),单位:毫秒
	field, err = w.CreateFormField("duration")
	if err != nil {
		return "", fmt.Errorf("create field error, got %v", err)
	}
	field.Write([]byte("3000"))

	fileWriter, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create file writer error, got %v", err)
	}
	io.Copy(fileWriter, fileReader)

	req, err := http.NewRequest("POST", bot.AddrForUploadFile(), postBody)
	if err != nil {
		return "", fmt.Errorf("failed to construct request, got %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.token))

	res, err := bot.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("upload file failed, status: %d", res.StatusCode)
	}

	type ResponseBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			FileKey string `json:"file_key"`
		} `json:"data"`
	}

	r := &ResponseBody{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return "", fmt.Errorf("can decode res body, err: %v", err)
	}

	return r.Data.FileKey, nil
}

func (bot *FeishuGroupBot) DownloadFile(fileKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", bot.AddrForDownloadFile(fileKey), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.token))

	res, err := bot.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("download file failed, status: %d", res.StatusCode)
	}

	result := []byte{}
	if _, err := res.Body.Read(result); err != nil {
		return nil, fmt.Errorf("read file body failed")
	}

	return result, nil
}
