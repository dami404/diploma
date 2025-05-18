package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dami404/diploma-db/config"
	"github.com/dami404/diploma-db/iternal/app/psql"
	"github.com/dami404/diploma-db/iternal/app/server"
	"github.com/dami404/diploma-db/iternal/usecase"

	"github.com/dami404/diploma-db/iternal/controller"
	"github.com/dami404/diploma-db/iternal/repository"
)

func Run(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg := psql.NewConnection(ctx, cfg.DBUrl)
	defer func() {
		if err := pg.Close(ctx); err != nil {
			log.Fatal(fmt.Sprintf("failed to close database connection: %v", err))
		} else {
			log.Println("database connection closed")
		}
	}()

	var err error
	dbRepo := repository.NewDBRepository(ctx, pg)
	dbUsecase := usecase.NewDBUsecase(dbRepo)
	dbhandler := controller.NewDBHandler(dbUsecase)

	log.Println("starting http server")
	httpServer := server.NewHttpServer(cfg.Host, cfg.Port)
	httpServer.Router = controller.SetRouter(dbhandler.LastEvents, dbhandler.Save)

	httpServer.Start()
	log.Println(fmt.Sprintf("listen and serve: %s:%s", cfg.Host, cfg.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err = <-httpServer.Notify():
		log.Fatal(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
	}

}
