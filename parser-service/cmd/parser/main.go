package main

import (
	"log"

	"github.com/dami404/diploma-parser/config"
	app "github.com/dami404/diploma-parser/iternal/app"
)

func main() {
	log.Println("reading config")
	cfg := config.MustLoad()
	log.Println("config loaded")
	log.Println("starting app")
	app.Run(cfg)
}
