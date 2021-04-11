package controller

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func (ctrl *Controller) CompileTemplate(text string) (*template.Template, error) {
	return template.New("_").Funcs(template.FuncMap{}).Funcs(sprig.TxtFuncMap()).Parse(text)
}

func (ctrl *Controller) RenderTemplate(tpl *template.Template, data interface{}) (string, error) {
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("render a template with params: %w", err)
	}
	return buf.String(), nil
}
