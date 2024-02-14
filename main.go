package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tn"
	app.Usage = "A simple terminal application for writing notes."
	app.Action = func(c *cli.Context) error {
		println("Hello from Terminal Doodle app!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
