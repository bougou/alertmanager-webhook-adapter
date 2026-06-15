package card

type Element interface {
	element()
}

type ActionElement interface {
	actionElement()
}

type ElemTextOrImage struct {
	*Text
	*ElemImage
}

// ElemImage 图像元素，注与图像模块的区别（ElementImg)
type ElemImage struct {
	Tag     string `json:"tag"` // tag=img
	ImgKey  string `json:"img_key"`
	Alt     *Text  `json:"alt"`
	Preview bool   `json:"preview,omitempty"`
}

func (e *ElemImage) element() {
}

type ElemButton struct {
	Tag      string                 `json:"tag"`                 // tag=button
	Text     *Text                  `json:"text"`                // 按钮中的文本
	URL      string                 `json:"url,omitempty"`       // 跳转链接，和multi_url互斥
	MultiURL *MultiURL              `json:"multi_url,omitempty"` // 多端跳转链接
	Type     string                 `json:"type,omitempty"`      // 配置按钮样式，默认为"default", "default"/"primary"/"danger"
	Value    map[string]interface{} `json:"value,omitempty"`     // 点击后返回业务方
	Confirm  *Confirm               `json:"confirm,omitempty"`   // 二次确认的弹框
}

func (e *ElemButton) element() {
}

func (e *ElemButton) actionElement() {
}

type ElemSelectMenu struct {
	Tag           string                 `json:"tag"` // "select_static" / "select_person", 元素标签，选项模式/选人模式
	PlaceHolder   *Text                  `json:"place_holder,omitempty"`
	InitialOption string                 `json:"initial_option,omitempty"`
	Options       []*Option              `json:"option,omitempty"`
	Value         map[string]interface{} `json:"value,omitempty"`
	Confirm       *Confirm               `json:"confirm,omitempty"`
}

func (e *ElemSelectMenu) element() {
}

func (e *ElemSelectMenu) actionElement() {
}

type ElemOverflow struct {
	Tag     string                 `json:"tag"` // tag=overflow
	Options []*Option              `json:"option"`
	Value   map[string]interface{} `json:"value,omitempty"`
	Confirm *Confirm               `json:"confirm,omitempty"`
}

func (e *ElemOverflow) element() {
}

func (e *ElemOverflow) actionElement() {
}

type ElemDatePicker struct {
	Tag             string                 `json:"tag"` // 如下三种取值 "date_picker", "picker_time", "picker_datetime"
	InitialDate     string                 `json:"initial_date,omitempty"`
	InitialTime     string                 `json:"initial_time,omitempty"`
	InitialDateTime string                 `json:"initial_datetime,omitempty"`
	PlaceHolder     *Text                  `json:"placeholder,omitempty"`
	Value           map[string]interface{} `json:"value,omitempty"`
	Confirm         *Confirm               `json:"confirm,omitempty"`
}

func (e *ElemDatePicker) element() {
}

func (e *ElemDatePicker) actionElement() {
}
