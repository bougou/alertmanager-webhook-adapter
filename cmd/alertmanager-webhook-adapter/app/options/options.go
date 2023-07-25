package options

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	restful "github.com/emicklei/go-restful/v3"

	"github.com/bougou/alertmanager-webhook-adapter/pkg/api"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/models"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/models/templates"
)

type AppOptions struct {
	Addr        string
	Signature   string
	TmplDir     string
	TmplName    string
	TmplDefault string
	TmplLang    string
	Version     bool
	Debug       bool
}

func NewAppOptions() *AppOptions {
	return &AppOptions{}
}

func (o *AppOptions) Run() error {
	execFile, err := os.Executable()
	if err != nil {
		panic("fatal")
	}

	// If using builtin templates (o.TmplDir == ""), then we must check whether or not the specified lang is supported.
	if o.TmplLang != "" && o.TmplDir == "" {
		if _, exists := templates.DefaultTmplByLang[o.TmplLang]; !exists {
			return fmt.Errorf("the builtin templates does not support specified lang: (%s)", o.TmplLang)
		}
		if err := models.LoadDefaultTemplate(o.TmplLang); err != nil {
			return fmt.Errorf("load default template for lang (%s) failed, err: %s", o.TmplLang, err)
		}
	}

	if o.TmplDir == "" && (o.TmplName != "" || o.TmplDefault != "") {
		fmt.Println("Warning, there is no meaning to specify --tmpl-name or --tmpl-default option without specify --tmpl-dir option, just ingored.")
	}

	if o.TmplDir != "" {
		if o.TmplName != "" && o.TmplDefault != "" {
			fmt.Println("Warning, there is no meaning to specify --tmpl-name and --tmpl-default options together, --tmpl-default is ignored.")
			o.TmplDefault = ""
		}

		if !filepath.IsAbs(o.TmplDir) {
			o.TmplDir = filepath.Join(filepath.Dir(execFile), o.TmplDir)
		}

		if err := models.LoadTemplate(o.TmplDir, o.TmplName, o.TmplDefault, o.TmplLang); err != nil {
			msg := fmt.Sprintf("Load templates from dir (%s) failed, err: %s", o.TmplDir, err)
			return errors.New(msg)
		}
	}

	container := restful.DefaultContainer

	controller := api.NewController(o.Signature)
	controller.WithDebug(o.Debug)

	controller.Install(container)

	s := &http.Server{
		Addr:    o.Addr,
		Handler: container,
	}
	log.Printf("start listening at %s", s.Addr)
	return s.ListenAndServe()
}
