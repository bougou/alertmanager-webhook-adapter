package models

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

var (
	topLevelTemplateName = "prom"

	// store the default templates
	promMsgTemplate *safeTemplate

	// store templates for different channels
	promMsgTemplatesMap = make(map[string]*safeTemplate)

	defaultFuncs = map[string]interface{}{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"title":   strings.Title,

		"markdown": markdownEscapeString,
	}
	isMarkdownSpecial [128]bool
)

func init() {
	var err error

	if err = LoadDefaultTemplate(templates.DefaultTmpl); err != nil {
		panic(err)
	}

	for _, c := range "_*`" {
		isMarkdownSpecial[c] = true
	}
}

func LoadDefaultTemplate(defaultTmpl string) error {
	promMsgTemplate = &safeTemplate{}
	if err := promMsgTemplate.UpdateTemplate(defaultTmpl); err != nil {
		msg := fmt.Sprintf("UpdateTemplate for default failed, err: %s", err)
		return errors.New(msg)
	}

	for k, v := range templates.ChannelsDefaultTmplMap {
		t := &safeTemplate{}
		if err := t.UpdateTemplate(v); err != nil {
			msg := fmt.Sprintf("UpdateTemplate for (%s) failed, err: %s", k, err)
			return errors.New(msg)
		}
		promMsgTemplatesMap[k] = t
	}

	return nil
}

// LoadTemplate loads external templates from specified template dir
func LoadTemplate(tmplDir, tmplName, tmplDefault string) error {

	// If tmplName is not empty, use the specified tmpl to update the default promMsgTemplate
	// and clear the promMsgTemplatesMap, thus will use the specified tmpl for all notification channels.
	if tmplName != "" {
		for k := range promMsgTemplatesMap {
			delete(promMsgTemplatesMap, k)
		}

		tmplFile := path.Join(tmplDir, fmt.Sprintf("%s.%s", tmplName, "tmpl"))
		b, err := os.ReadFile(tmplFile)
		if err != nil {
			msg := fmt.Sprintf("read file (%s) failed, err: %s", tmplFile, err)
			return errors.New(msg)
		}

		if err := promMsgTemplate.UpdateTemplate(string(b)); err != nil {
			msg := fmt.Sprintf("UpdateTemplate for default failed, err: %s", err)
			return errors.New(msg)
		}

		return nil
	}

	var customDefaultTmpl string
	if tmplDefault != "" {
		tmplFile := path.Join(tmplDir, fmt.Sprintf("%s.%s", tmplDefault, "tmpl"))
		b, err := os.ReadFile(tmplFile)
		if err != nil {
			msg := fmt.Sprintf("read file (%s) failed, err: %s", tmplFile, err)
			return errors.New(msg)
		}
		customDefaultTmpl = string(b)
	}

	// try to find template file named "<channel>.tmpl" and update the promTemplatesMap
	for channel, t := range promMsgTemplatesMap {
		var channelTmpl string

		tmplFile := path.Join(tmplDir, fmt.Sprintf("%s.%s", channel, "tmpl"))
		b, err := os.ReadFile(tmplFile)

		if os.IsNotExist(err) {
			// case 1: <channel>.tmpl file does not exist, and not specified custom default, use continue the next loop
			if tmplDefault == "" {
				continue
			}
			// case 2: <channel>.tmpl file does not exist, but specified custom default, use custom default as tmpl
			channelTmpl = customDefaultTmpl
		} else {
			// case 3: <channel>.tmpl exists, but read failed, error and return
			if err != nil {
				msg := fmt.Sprintf("read file (%s) failed, err: %s", tmplFile, err)
				return errors.New(msg)
			}
			// case 4: <channel>.tmpl exists, and read succeeded, use file content as tmpl
			channelTmpl = string(b)
		}

		if err := t.UpdateTemplate(channelTmpl); err != nil {
			msg := fmt.Sprintf("UpdateTemplate for (%s) failed, err: %s", channel, err)
			return errors.New(msg)
		}
	}

	return nil
}

type safeTemplate struct {
	*template.Template
	current string
	mu      sync.RWMutex
}

func (t *safeTemplate) UpdateTemplate(newTpl string) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tpl, err := template.New(topLevelTemplateName).
		Funcs(defaultFuncs).
		Option("missingkey=zero").
		Parse(newTpl)
	if err != nil {
		return
	}

	_ = t.current // oldtmpl
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

	tmpl, err := promMsgTemplate.Clone()
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
