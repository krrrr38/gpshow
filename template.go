package main

import (
	"bytes"
	"text/template"

	"github.com/krrrr38/gpshow/utils"
)

// SlideContent contain template parameters
type SlideContent struct {
	Title  string
	Slides string

	Asset     string
	JQuery    string
	Prettyfy  string
	Languages []string
}

// Render picture-show index page html
func Render(content *SlideContent) []byte {
	tpl, err := Asset("resources/slide.tpl")
	utils.DieIf(err)

	var buf bytes.Buffer
	tmpl := template.Must(template.New("slide").Parse(string(tpl[:])))
	err = tmpl.Execute(&buf, content)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
