package card

// Text 文本对象
type Text struct {
	// Tag 支持"plain_text"和"lark_md"两种模式
	// https://open.feishu.cn/document/ukTMukTMukTM/uADOwUjLwgDM14CM4ATN
	Tag     string `json:"tag"`
	Content string `json:"content"`         // 文本内容
	Lines   int    `json:"lines,omitempty"` // 内容显示行数, 1 显示行数， lines字段仅支持"plain_text"模式
	I18N    *I18N  `json:"i18n,omitempty"`
}

type I18N struct {
	ZHCN string `json:"zh_cn"`
	ENUS string `json:"en_us"`
	JAJP string `json:"jn_jp"`
}

type Field struct {
	IsShort bool  `json:"is_short"`
	Text    *Text `json:"text"`
}

type MultiURL struct {
	URL        string `json:"url"`
	AndroidURL string `json:"android_url"`
	IOSURL     string `json:"ios_url"`
	PCURL      string `json:"pc_url"`
}

type Option struct {
	Value    string    `json:"value"`          // 选项选中后返回业务方的数据
	Text     *Text     `json:"text,omitempty"` // 选项显示内容，非待选人员时必填
	URL      string    `json:"url,omitempty"`  // *仅支持overflow，跳转指定链接，和multi_url字段互斥
	MultiURL *MultiURL `json:"multiURL,omitempty"`
}

type Confirm struct {
	Title *Text `json:"title"` // 弹框标题, 仅支持"plain_text"
	Text  *Text `json:"text"`  // 弹框内容, 仅支持"plain_text"
}
