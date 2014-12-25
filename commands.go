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
	err := os.Mkdir(projectName, 0777)
	utils.DieIf(err)

	templateDir := "resources/project_template"
	templateFiles, err := AssetDir(templateDir)
	utils.DieIf(err)

	for _, filename := range templateFiles {
		generateProjectTemplates(projectName, templateDir, filename)
	}
	utils.Log("info", fmt.Sprintf("Generate new slide project directory: %s", projectName))
}

func generateProjectTemplates(projectDir, templateDir, filename string) {
	templatePath := templateDir + "/" + filename
	projectPath := projectDir + "/" + filename
	_, err := AssetInfo(templatePath)

	if err != nil {
		err := os.Mkdir(projectPath, 0777)
		utils.DieIf(err)

		templateFiles, err := AssetDir(templatePath)
		utils.DieIf(err)
		for _, nextFilename := range templateFiles {
			generateProjectTemplates(projectPath, templatePath, nextFilename)
		}
	} else {
		bytes, err := Asset(templatePath)
		utils.DieIf(err)
		ioutil.WriteFile(projectPath, bytes, os.ModePerm)
	}
}

func doOffline(c *cli.Context) {
	// outDir := c.String("output")
	showPath := c.String("show")
	showPathInfo, err := os.Stat(showPath)
	utils.DieIf(err)
	if !showPathInfo.IsDir() {
		utils.Log("error", fmt.Sprintf("The path `%s` is not an accessible directory", showPath))
		os.Exit(1)
	}

	config := ConfigFile(showPath + "/conf.js")
	adapter := &OfflineAdapter{
		showPath: showPath,
		config:   config,
	}
	html := adapter.HTML()
	utils.Log("warn", string(html))
	// TODO
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
