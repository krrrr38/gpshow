package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/swdyh/go-enumerable/src/enumerable"
	"regexp"
	"sync"
)

// MakeSlide generate html from markdown text
func MakeSlide(pageTotal *int, markdown []byte) []byte {
	var buffer bytes.Buffer
	rawPages := regexp.MustCompile("(?m)^!SLIDE").Split(string(markdown[:]), -1)

	// remove empty page
	var nonContentFilter func([]string) []string
	enumerable.MakeFilter(&nonContentFilter, func(page string) bool { return page != "" })
	pages := nonContentFilter(rawPages)

	htmls := make([]string, len(pages))
	var wg sync.WaitGroup
	for i, page := range pages {
		wg.Add(1)
		go func(i int, page string) {
			content := blackfriday.MarkdownBasic([]byte(page))
			htmls[i] = fmt.Sprintf("<div class=\"content\" id=\"slide-%d\"><div class=\"container\">%s</div></div>", i+(*pageTotal), string(content[:]))
			wg.Done()
		}(i, page)
	}
	wg.Wait()
	(*pageTotal) += len(pages)

	for _, html := range htmls {
		buffer.WriteString(html)
	}
	return buffer.Bytes()
}
