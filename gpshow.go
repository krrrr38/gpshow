package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gpshow"
	app.Version = Version
	app.Usage = "slip and slide picture shows for the web\n\tGo port of picture-show (https://github.com/softprops/picture-show)"
	app.Author = "krrrr38"
	app.Email = "k.kaizu38@gmail.com"

	app.Action = PShowAction
	app.Flags = PShowFlags
	app.Commands = Commands

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	app.Run(os.Args)
}
