package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dami404/diploma-parser/config"
	"github.com/dami404/diploma-parser/iternal/app/server"
	"github.com/dami404/diploma-parser/iternal/controller"
	"github.com/dami404/diploma-parser/iternal/repository"
	"github.com/dami404/diploma-parser/iternal/usecase"
)

func Run(cfg *config.Config) {
	dbRepo := repository.NewDBRepository(cfg.DBUrl)
	parserUsecase := usecase.NewParserUsecase(dbRepo)
	handler := controller.NewParserHandler(parserUsecase)

	log.Println("starting http server")
	httpServer := server.NewHttpServer(cfg.Host, cfg.Port)
	httpServer.Router = controller.SetRouter(handler.SearchEvents)

	httpServer.Start()
	log.Println(fmt.Sprintf("listen and serve: %s:%s", cfg.Host, cfg.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err := <-httpServer.Notify():
		log.Fatal(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
	}
}
