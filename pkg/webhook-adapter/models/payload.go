package models

// The purpse of Payload is to hide the complexity of constructing channel-specific Msg.
// Because each specific channel provides suitable Payload2MsgFn(s) convertion functions for its supported msgType(s).
type Payload struct {
	Raw      string   `json:"raw"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Markdown string   `json:"markdown"` // Don't put Title content in Markdown
	Files    []string `json:"files"`
	Images   []Image  `json:"images"`
	Links    []Link   `json:"links"`
	Buttons  []Button `json:"buttons"`
	At       At       `json:"at"`
}

type PayloadGenerator interface {
	ToPayload() *Payload
}

type Image struct {
	Bytes  []byte `json:"bytes"`
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

type Link struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	PicURL string `json:"picURL"`
}

type Button struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtAll     bool     `json:"atAll"`
}
