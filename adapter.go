package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/github"
	"github.com/krrrr38/gpshow/utils"
)

// assets path for static contents of server
const AssetsPath = "/assets/"
const prettify = "http://cdnjs.cloudflare.com/ajax/libs/prettify/r298/"
const jquery = "http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"

var languages = []string{"apollo", "basic", "clj", "css", "dart", "erlang", "go", "hs", "lisp", "llvm", "lua", "matlab", "ml", "mumps", "n", "pascal", "proto", "r", "rd", "scala", "sql", "tcl", "tex", "vb", "vhdl", "wiki", "xq", "yaml"}

// SlideMaker can generate slide component
type SlideMaker interface {
	HTML() []byte
}

// DefaultAdapter contain local picture-show project info
type DefaultAdapter struct {
	showPath string
	config   Configuration
}

// OfflineAdapter contain local picture-show project info
type OfflineAdapter struct {
	showPath string
	config   Configuration
}

// GistAdapter contain gist info
type GistAdapter struct {
	id string
}

// HTML generate slide html based on config in base dir
func (adapter *DefaultAdapter) HTML() []byte {
	var buffer bytes.Buffer
	showPath := adapter.showPath
	var pageTotal = 0
	for _, section := range adapter.config.Sections {
		markdownPath := fmt.Sprintf("%s/%s/%s.md", showPath, section, section)
		markdown, err := ioutil.ReadFile(markdownPath)
		if err != nil {
			utils.Log("warn", fmt.Sprintf("no file(s) at path %s/%s/%s.md", showPath, section, section))
		}
		buffer.Write(MakeSlide(&pageTotal, markdown))
	}
	content := &SlideContent{
		Title:     adapter.config.Title,
		Slides:    buffer.String(),
		Asset:     AssetsPath,
		JQuery:    jquery,
		Prettyfy:  prettify,
		Languages: languages,
	}
	return Render(content)
}

// HTML generate static slide html based on config in base dir
func (adapter *OfflineAdapter) HTML() []byte {
	return nil // TODO
}

// HTML generate slide html based on config in base dir
func (adapter *GistAdapter) HTML() []byte {
	var buffer bytes.Buffer
	var pageTotal = 0
	var title string
	gist := FetchGist(adapter.id)
	files := gist.Files
	if conffile, ok := files["conf.js"]; ok {
		config := Config([]byte(*conffile.Content))
		title = config.Title
		for _, section := range config.Sections {
			if markdown, ok := files[github.GistFilename(fmt.Sprintf("%s.md", section))]; ok {
				buffer.Write(MakeSlide(&pageTotal, []byte(*markdown.Content)))
			} else {
				utils.Log("warn", fmt.Sprintf("cannot find %s.md in this gist", section))
			}
		}
	} else {
		title = *gist.Description
		for _, file := range files {
			if strings.HasSuffix(*file.Filename, ".md") {
				buffer.Write(MakeSlide(&pageTotal, []byte(*file.Content)))
			} else {
				utils.Log("warn", fmt.Sprintf("`%s` is not have `.md` ext.", *file.Filename))
			}
		}
	}

	content := &SlideContent{
		Title:     title,
		Slides:    buffer.String(),
		Asset:     AssetsPath,
		JQuery:    jquery,
		Prettyfy:  prettify,
		Languages: languages,
	}
	return Render(content)
}
