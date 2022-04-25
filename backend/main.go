package main

import (
	"github.com/Qwiri/InnoDays2022/backend/internal"
	"github.com/Qwiri/InnoDays2022/backend/internal/server"
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func main() {

	// load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Warn("Could not load .env file")
	}

	// connect to db
	var db *gorm.DB
	if os.Getenv("prod") != "prod" {
		_db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
		if err != nil {
			log.WithError(err).Fatal("Could not connect to sqlite db")
			return
		}
		db = _db
	} else {
		dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Europe/Berlin"
		_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.WithError(err).Fatal("Could not connect to postgres db")
			return
		}
		db = _db
	}
	log.Info("Connected to database")

	// automigrate db
	if err := db.AutoMigrate(internal.TableModels...); err != nil {
		log.WithError(err).Fatal("Could not migrate db")
		return
	}
	log.Info("Migrated database")

	// initialize server
	var s server.Server
	s.New(db)
	s.ConnectRoutes()
	if err := s.Listen(":3001"); err != nil {
		log.WithError(err).Fatal("Could not start server")
	}

}
