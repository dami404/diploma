package main

import (
	"github.com/dami404/diploma-db/config"
	"github.com/dami404/diploma-db/iternal/app"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
