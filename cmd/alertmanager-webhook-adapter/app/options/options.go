package options

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
			return fmt.Errorf("the builtin templates does not support specified lang (%s), builtin supported langs: (%s)", o.TmplLang, strings.Join(templates.DefaultSupportedLangs(), ","))
		}
		if err := models.LoadDefaultTemplate(o.TmplLang); err != nil {
			return fmt.Errorf("load default template for lang (%s) failed, err: %s", o.TmplLang, err)
		}
	}

	if o.TmplDir == "" && (o.TmplName != "" || o.TmplDefault != "") {
		fmt.Println("Warning, there is no meaning to specify --tmpl-name or --tmpl-default option without specify --tmpl-dir option, just ignored.")
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

	fmt.Println("Signature: ", o.Signature)
	if o.Signature == "未知" {
		fmt.Println("Warn, you are using the default signature, we suggest to specify a custom signature by --signature option.")
	}

	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")
	noProxy := os.Getenv("NO_PROXY")
	if httpProxy != "" || httpsProxy != "" {
		fmt.Println("Found http proxy from environment variables:")
		fmt.Printf("  HTTP_PROXY: (%s)\n", httpProxy)
		fmt.Printf("  HTTPS_PROXY: (%s)\n", httpsProxy)
		fmt.Printf("  NO_PROXY: (%s)\n", noProxy)
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
