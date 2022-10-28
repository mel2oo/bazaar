package main

import (
	"bazaar/config"
	"bazaar/internal/app"
	"flag"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "../../config/config.toml", "config file")
	flag.Parse()

	c, err := config.Load(path)
	if err != nil {
		panic(err)
	}

	if err := app.Run(c); err != nil {
		panic(err)
	}
}
