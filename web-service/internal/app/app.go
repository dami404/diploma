package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dami404/diploma-web/config"
	"github.com/dami404/diploma-web/internal/app/server"
	"github.com/dami404/diploma-web/internal/controller"
	"github.com/dami404/diploma-web/internal/repository"
	"github.com/dami404/diploma-web/internal/usecase"
)

func Run(cfg *config.Config) {

	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	// defer cancel()

	appRepo := repository.NewHTTPRepository(cfg.DBUrl, cfg.Url)
	appUsecase := usecase.NewUsecase(appRepo)
	handler := controller.NewHandler(appUsecase)

	log.Println("starting http server")
	httpServer := server.NewHttpServer(cfg.Host, cfg.Port)
	httpServer.Router = controller.SetRouter(handler.StartPage, handler.Results, handler.LastEvents)

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
