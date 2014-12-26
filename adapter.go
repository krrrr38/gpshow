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
const AssetsPath = "assets/"
const prettifyURLBase = "http://cdnjs.cloudflare.com/ajax/libs/prettify/r298/"
const jqueryURL = "http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"

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
	outDir   string
}

// GistAdapter contain gist info
type GistAdapter struct {
	id string
}

// HTML generate slide html based on config in base dir
func (adapter *DefaultAdapter) HTML() []byte {
	slides := generateLocalSlides(adapter.showPath, adapter.config)
	content := &SlideContent{
		Title:     adapter.config.Title,
		Slides:    slides,
		Asset:     "/" + AssetsPath,
		JQuery:    jqueryURL,
		Prettyfy:  prettifyURLBase,
		Languages: languages,
	}
	return Render(content)
}

// HTML generate static slide html based on config in base dir
func (adapter *OfflineAdapter) HTML() []byte {
	showPath := adapter.showPath
	slides := generateLocalSlides(showPath, adapter.config)

	CopyResourceDir("assets", adapter.outDir+"/assets")
	CopyLocalStaticFiles(adapter.outDir, showPath, adapter.config.Sections)

	prettifyBasePath := AssetsPath + "prettify/"
	downloadFiles := []StaticDownload{
		StaticDownload{jqueryURL, AssetsPath + "js/", "jquery.min.js"},
		StaticDownload{prettifyURLBase + "prettify.css", prettifyBasePath, "prettify.css"},
		StaticDownload{prettifyURLBase + "prettify.js", prettifyBasePath, "prettify.js"},
	}
	for _, lang := range languages {
		filename := fmt.Sprintf("lang-%s.js", lang)
		langFile := StaticDownload{prettifyURLBase + filename, prettifyBasePath, filename}
		downloadFiles = append(downloadFiles, langFile)
	}
	DownloadStaticFiles(adapter.outDir, downloadFiles)

	content := &SlideContent{
		Title:     adapter.config.Title,
		Slides:    slides,
		Asset:     AssetsPath,
		JQuery:    AssetsPath + "js/jquery.min.js",
		Prettyfy:  prettifyBasePath,
		Languages: languages,
	}
	return Render(content)
}

func generateLocalSlides(showPath string, config Configuration) string {
	var buffer bytes.Buffer
	var pageTotal = 0
	for _, section := range config.Sections {
		markdownPath := fmt.Sprintf("%s/%s/%s.md", showPath, section, section)
		markdown, err := ioutil.ReadFile(markdownPath)
		if err != nil {
			utils.Log("warn", fmt.Sprintf("no file(s) at path %s/%s/%s.md", showPath, section, section))
		}
		buffer.Write(MakeSlide(&pageTotal, markdown))
	}
	return buffer.String()
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
		JQuery:    jqueryURL,
		Prettyfy:  prettifyURLBase,
		Languages: languages,
	}
	return Render(content)
}
