package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/darchlabs/backoffice/config"

	"github.com/darchlabs/backoffice/internal/api"
	"github.com/darchlabs/backoffice/internal/application"
	"github.com/darchlabs/backoffice/internal/storage"
	_ "github.com/darchlabs/backoffice/migrations"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	log.Println("Starting B-Cards backoffice")
	conf := &config.Config{}
	err := envconfig.Process("", conf)
	check(err)

	log.Printf("Database postgresql connection [B-Cards backoffice]\n")
	sqlStore, err := storage.NewSQLStore(conf.DBDriver, conf.DBDsn)
	check(err)

	err = goose.Up(sqlStore.DB.DB, "migrations/")
	check(err)

	log.Println("Starting B-Cards backoffice")
	app, err := application.New(&application.Config{
		Config:   conf,
		SqlStore: sqlStore,
	})
	check(err)

	server := api.NewServer(&api.ServerConfig{
		Port: conf.ApiServerPort,
		App:  app,
	})

	log.Printf("Starting [B-Cards backoffice]\n")
	err = server.Start(app)
	check(err)

	// listen interrupt
	quit := make(chan struct{})
	listenInterruption(quit)
	<-quit
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func listenInterruption(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("signal received", s.String())
		quit <- struct{}{}
	}()
}
