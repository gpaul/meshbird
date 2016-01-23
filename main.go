package main

import (
	"os"

	"github.com/gophergala2016/meshbird/core"
	"github.com/codegangsta/cli"
)

var (
// VERSION var using for auto versioning through Go linker
	VERSION = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "meshbird"
	app.Usage = "distributed overlay private networking"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "create new network",
			Action:  actionNew,
		},
		{
			Name:    "join",
			Aliases: []string{"r"},
			Usage:   "join network",
			Action:  actionJoin,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		println("error: ", err)
	}
}

func actionNew(ctx *cli.Context) {

}

func actionJoin(ctx *cli.Context) {
	node := core.NewNode(nil)
	defer node.Stop()
}