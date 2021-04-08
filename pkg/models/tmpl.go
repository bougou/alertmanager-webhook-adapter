package models

import (
	"bytes"
	_ "embed"
	"strings"
	"sync"
	"text/template"
)

//go:embed default.tmpl
var defaultTmpl string

var (
	topLevelTemplateName = "prom"
	promTemplate         safeTemplate
	defaultFuncs         = map[string]interface{}{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"title":   strings.Title,
		// join is equal to strings.Join but inverts the argument order
		// for easier pipelining in templates.
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"markdown": markdownEscapeString,
	}
	isMarkdownSpecial [128]bool
)

func init() {
	var err error

	_, err = LoadDefaultTemplate(defaultTmpl)
	if err != nil {
		panic(err)
	}

	for _, c := range "_*`" {
		isMarkdownSpecial[c] = true
	}
}

func LoadDefaultTemplate(newTpl string) (oldTpl string, err error) {
	return promTemplate.UpdateTemplate(newTpl)
}

type safeTemplate struct {
	*template.Template
	current string
	mu      sync.RWMutex
}

func (t *safeTemplate) UpdateTemplate(newTpl string) (oldTpl string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tpl, err := template.New(topLevelTemplateName).
		Funcs(defaultFuncs).
		Option("missingkey=zero").
		Parse(newTpl)
	if err != nil {
		return
	}

	oldTpl = t.current
	t.Template = tpl
	t.current = newTpl
	return
}

func (t *safeTemplate) Clone() (*template.Template, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.Template.Clone()
}

func markdownEscapeString(s string) string {
	b := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(b)

	for _, c := range s {
		if c < 128 && isMarkdownSpecial[c] {
			buf.WriteByte('\\')
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func ExecuteTextString(text string, data interface{}) (string, error) {
	if text == "" {
		return "", nil
	}

	tmpl, err := promTemplate.Clone()
	if err != nil {
		return "", err
	}

	tmpl, err = tmpl.New("").Option("missingkey=zero").Parse(text)
	if err != nil {
		return "", err
	}

	// reserve a buffer in 1k
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
