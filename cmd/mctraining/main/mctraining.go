package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/materials-commons/mctraining/cmd/mctraining"
)

func init() {
	mctraining.DBConnect()
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "V. Glenn Tarcea",
			Email: "gtarcea@umich.edu",
		},
	}
	app.Commands = []cli.Command{
		mctraining.CreateCommand,
	}

	app.Run(os.Args)
}
