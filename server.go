package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/krrrr38/gpshow/utils"
)

type staticBinaryHandler struct{}

// Server for picture-show slides
func Server(port int, slidemaker SlideMaker) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(slidemaker.HTML())
	})
	for _, path := range []string{"css", "js", "images"} {
		dir := fmt.Sprintf("/%s/", path)
		http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir(path))))
	}
	http.Handle(AssetsPath, &staticBinaryHandler{})

	utils.Log("info", fmt.Sprintf("starting show on http://localhost:%d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (h *staticBinaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	bytes, err := Asset("resources" + path)
	if err == nil {
		if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}
