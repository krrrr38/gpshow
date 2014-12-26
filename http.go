package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/krrrr38/gpshow/utils"
	"io/ioutil"
	"net/http"
)

// FetchGist just get gist data through github api
func FetchGist(id string) *github.Gist {
	client := github.NewClient(nil)
	gist, res, err := client.Gists.Get(id)
	defer res.Body.Close()
	utils.DieIf(err)

	return gist
}

// FetchFile just send GET request
func FetchFile(url string) (bytes []byte, err error) {
	res, err := http.Get(url)
	defer res.Body.Close()
	if res.StatusCode == 200 {
		bytes, err = ioutil.ReadAll(res.Body)
	} else {
		err = fmt.Errorf("Cannot find file %s : %d", url, res.StatusCode)
	}
	return bytes, err
}
