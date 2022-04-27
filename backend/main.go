package main

import (
	"context"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/Qwiri/InnoDays2022/backend/internal/server"
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Warn("Could not load .env file")
	}

	// connect to db
	var (
		db  *gorm.DB
		err error
	)
	if os.Getenv("prod") != "prod" {
		if db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{}); err != nil {
			log.WithError(err).Fatal("Could not connect to sqlite db")
			return
		}
	} else {
		dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Europe/Berlin"
		if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			log.WithError(err).Fatal("Could not connect to postgres db")
			return
		}
	}
	log.Info("Connected to database")

	// auto migrate db
	if err := db.AutoMigrate(common.TableModels...); err != nil {
		log.WithError(err).Fatal("Could not migrate db")
		return
	}
	log.Info("Migrated database")

	// initialize server
	s := server.New(db)

	// start web server
	go func() {
		if err := s.Listen(":3000"); err != nil {
			log.WithError(err).Fatal("Could not start server")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start janitor
	j := server.NewJanitor(db, s, ctx, time.Second*10)
	go j.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	cancel()
	log.Info("Shutting down")
	if err := s.Shutdown(); err != nil {
		log.WithError(err).Warn("cannot shutdown web server")
	}
	log.Info("Bye Bye!")
}
