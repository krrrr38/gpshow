package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/krrrr38/gpshow/utils"
)

// gpshow sub command list
var Commands = []cli.Command{
	commandInit,
	commandOffline,
	commandGist,
}

// gpshow main command flags
var PShowFlags = []cli.Flag{
	cli.IntFlag{Name: "port, p", Value: 3000, Usage: "port number"},
	cli.StringFlag{Name: "show, s", Value: ".", Usage: "path to show", EnvVar: "SHOW_HOME"},
}

var commandInit = cli.Command{
	Name:  "init",
	Usage: "Create picture-show project",
	Description: `
    Create picture-show project with <project_name>
`,
	Action: doInit,
}

var commandOffline = cli.Command{
	Name:  "offline",
	Usage: "Output static files based on current dir project",
	Description: `
    Output static files based on current dir project into output dir
`,
	Action: doOffline,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "show, s", Value: ".", Usage: "path to show", EnvVar: "SHOW_HOME"},
		cli.StringFlag{Name: "output, o", Value: "out", Usage: "Output directory"},
	},
}

var commandGist = cli.Command{
	Name:  "gist",
	Usage: "Start slides based on gist url which contains markdown texts",
	Description: `
    Start slides based on <gist_url> which contains markdown texts
`,
	Action: doGist,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "port, p", Value: 3000, Usage: "port number"},
	},
}

// PShowAction contains gpshow main action
func PShowAction(c *cli.Context) {
	port := c.Int("port")
	showPath := c.String("show")
	showPathInfo, err := os.Stat(showPath)
	utils.DieIf(err)
	if !showPathInfo.IsDir() {
		utils.Log("error", fmt.Sprintf("The path `%s` is not an accessible directory", showPath))
		os.Exit(1)
	}

	config := ConfigFile(showPath + "/conf.js")
	adapter := &DefaultAdapter{
		showPath: showPath,
		config:   config,
	}
	Server(port, adapter)
}

func doInit(c *cli.Context) {
	projectName := c.Args().First()
	if projectName == "" {
		cli.ShowCommandHelp(c, "init")
		os.Exit(1)
	}
	err := os.MkdirAll(projectName, 0755)
	utils.DieIf(err)

	CopyResourceDir("project_template", projectName)

	utils.Log("info", fmt.Sprintf("Generate new slide project directory: %s", projectName))
}

func doOffline(c *cli.Context) {
	outDir := c.String("output")
	showPath := c.String("show")
	showPathInfo, err := os.Stat(showPath)
	utils.DieIf(err)
	if !showPathInfo.IsDir() {
		utils.Log("error", fmt.Sprintf("The path `%s` is not an accessible directory", showPath))
		os.Exit(1)
	}

	config := ConfigFile(showPath + "/conf.js")

	os.RemoveAll(outDir)
	utils.DieIf(os.Mkdir(outDir, 0755))

	adapter := &OfflineAdapter{
		showPath: showPath,
		config:   config,
		outDir:   outDir,
	}
	html := adapter.HTML()

	ioutil.WriteFile(fmt.Sprintf("%s/index.html", outDir), html, 0644)
	utils.Log("info", fmt.Sprintf("create static slides in `%s` directory.", outDir))
}

func doGist(c *cli.Context) {
	port := c.Int("port")
	arg := c.Args().First()
	if arg == "" {
		cli.ShowCommandHelp(c, "gist")
		os.Exit(1)
	}

	// arg should be following cases
	// 1. https://gist.github.com/krrrr38/d709ca3a8cf92e4294b7
	// 2. https://gist.github.com/d709ca3a8cf92e4294b7.git
	// 3. git@gist.github.com:/d709ca3a8cf92e4294b7.git
	// 4. d709ca3a8cf92e4294b7
	// then, id = d709ca3a8cf92e4294b7
	id := extractGistID(arg)

	adapter := &GistAdapter{
		id: id,
	}
	Server(port, adapter)
}

func extractGistID(arg string) string {
	parts := strings.Split(arg, "/")
	return strings.Replace(parts[len(parts)-1], ".git", "", -1)
}
