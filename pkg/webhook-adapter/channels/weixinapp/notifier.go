package weixinapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ErrCodeAboutTokens contains weixinapp
// https://work.weixin.qq.com/api/doc/90000/90139/90313
var ErrCodeAboutTokens = []int{
	40014, // 不合法的access_token
	42001, // access_token已过期
}

// https://developer.work.weixin.qq.com/document/path/90236
type Notifier struct {
	addr        string
	corpID      string // 企业ID
	agentID     int    // 应用ID
	agentSecret string // 应用的凭证密钥, 区分应用
	client      *http.Client

	// toUser 指定接收消息的成员，成员 ID 列表（多个接收者用 '|' 分隔，最多支持 1000 个）。
	// 指定为 "@all"，则向该企业应用的全部成员发送
	toUser string
	// toParty 指定接收消息的部门，部门 ID 列表，多个接收者用 '|' 分隔，最多支持 100 个。
	// 当 touser 为 "@all" 时忽略本参数
	toParty string
	// toTag 指定接收消息的标签，标签ID列表，多个接收者用 '|' 分隔，最多支持 100 个。
	// 当 touser 为 "@all" 时忽略本参数
	toTag string

	token          string
	tokenAt        time.Time
	tokenExpiredIn time.Duration
}

// toUser,toParty,toTag CAN NOT be empty at the same time.
func NewNotifer(corpID string, agentID int, agentSecret string, toUser string, toParty string, toTag string) *Notifier {
	return &Notifier{
		addr:        "https://qyapi.weixin.qq.com",
		corpID:      corpID,
		agentID:     agentID,
		agentSecret: agentSecret,
		client:      &http.Client{},

		toUser:  toUser,
		toParty: toParty,
		toTag:   toTag,
	}
}

func (n *Notifier) Addr() string {
	return fmt.Sprintf("%s/cgi-bin/message/send?access_token=%s", n.addr, n.token)
}

func (n *Notifier) AddrForGetToken() string {
	return fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s", n.addr, n.corpID, n.agentSecret)
}

func (n *Notifier) GetToken() error {
	req, err := http.NewRequest("GET", n.AddrForGetToken(), nil)
	if err != nil {
		return fmt.Errorf("get token construct http request failed, %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("get token request failed, %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("get token response error, status: %d", res.StatusCode)
	}

	type ResponseBody struct {
		ErrCode          int    `json:"errcode"`
		ErrMsg           string `json:"errmsg"`
		AccessToken      string `json:"access_token"`
		ExpiresInSeconds int    `json:"expires_in"`
	}

	r := &ResponseBody{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return fmt.Errorf("get token decode response body failed, %s", err)
	}

	if r.ErrCode != 0 {
		return fmt.Errorf("get token failed, errmsg is %s", r.ErrMsg)
	}

	n.token = r.AccessToken
	n.tokenAt = time.Now()
	n.tokenExpiredIn, err = time.ParseDuration(fmt.Sprintf("%ds", r.ExpiresInSeconds))
	if err != nil {
		n.tokenExpiredIn = 2 * time.Hour
	}

	return nil

}

func (n *Notifier) ShouldGetToken() bool {
	if n.token == "" || time.Since(n.tokenAt) > n.tokenExpiredIn {
		return true
	}

	return false
}

func (n *Notifier) Send(msg *Msg) error {
	// fill agentID
	msg.AgentID = n.agentID
	msg.ToUser = n.toUser
	msg.ToParty = n.toParty
	msg.ToTag = n.toTag

	if err := validateMsg(msg.MsgType, msg); err != nil {
		return fmt.Errorf("valid msg failed, err: %s", err)
	}

	if n.token == "" {
		if err := n.GetToken(); err != nil {
			return fmt.Errorf("failed token failed, %s", err)
		}
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", n.Addr(), bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("failed to construct request, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send msg error, %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("send msg response error, status: %d", res.StatusCode)
	}

	return nil
}
