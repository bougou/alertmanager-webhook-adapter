package card

// ModuleDiv 内容模块
type ModuleDiv struct {
	Tag    string   `json:"tag"`             // div
	Text   *Text    `json:"text,omitempty"`  // 单个文本展示, 和 field 至少要有一个
	Fields []Field  `json:"field,omitempty"` // 多个文本展示, 和 text 至少要有一个
	Extra  *Element `json:"extra,omitempty"` // 展示附加元素, 最多可展示一个元素
}

func (e *ModuleDiv) cardModule() string {
	return "div"
}

// div, hr, img, note, action
type CardModule interface {
	cardModule() string
}

// ModuleHR 分割线模块
type ModuleHR struct {
	Tag string `json:"tag"` // hr
}

func (e *ModuleHR) cardModule() string {
	return "hr"
}

// ModuleImg 图片模块
type ModuleImg struct {
	Tag     string `json:"tag"` // img
	ImgKey  string `json:"img_key"`
	Title   *Text  `json:"title"`
	Mode    string `json:"mode,omitempty"` // 图片显示模式: crop_center：居中裁剪模式, fit_horizontal：平铺模式
	Alt     *Text  `json:"text,omitempty"`
	Preview bool   `json:"preview,omitempty"` // 点击后是否放大图片，缺省为true。在配置 card_link 后可设置为false，使用户点击卡片上的图片也能响应card_link链接跳转
}

func (e *ModuleImg) cardModule() string {
	return "img"
}

// ModuleNote 备注模块
type ModuleNote struct {
	Tag      string             `json:"tag"` // note
	Elements []*ElemTextOrImage `json:"elements"`
	// image
}

func (e *ModuleNote) cardModule() string {
	return "note"
}

// ModuleAction 交互模块
type ModuleAction struct {
	Tag     string          `json:"tag"`
	Actions []ActionElement `json:"actions"`
	Layout  string          `json:"layout,omitempty"` // bisected 为二等分布局, trisection 为三等分布局, flow 为流式布局元素会按自身大小横向排列并在空间不够的时候折行
}

func (e *ModuleAction) cardModule() string {
	return "action"
}

type ModuleMarkdown struct {
	Tag     string    `json:"tag"` // markdown
	Content string    `json:"content"`
	Href    *MultiURL `json:"href,omitempty"` // 差异化跳转：仅在需要PC、移动端跳转不同链接使用
}

func (e *ModuleMarkdown) cardModule() string {
	return "markdown"
}

func NewModuleMarkdown(content string, href *MultiURL) *ModuleMarkdown {
	return &ModuleMarkdown{
		Tag:     "markdown",
		Content: SanitizeMarkdown(content),
		Href:    href,
	}
}

func SanitizeMarkdown(content string) string {
	return content
}
