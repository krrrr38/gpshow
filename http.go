package main

import (
	"github.com/google/go-github/github"
	"github.com/krrrr38/gpshow/utils"
)

func FetchGist(id string) *github.Gist {
	client := github.NewClient(nil)
	gist, _, err := client.Gists.Get(id)
	utils.DieIf(err)

	return gist
}
