package main

import (
	"github.com/dami404/diploma-web/config"
	"github.com/dami404/diploma-web/internal/app"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
