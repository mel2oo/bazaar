package main

import (
	"bazaar/internal/cli"
	"os"
)

func main() {
	app := cli.New()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
