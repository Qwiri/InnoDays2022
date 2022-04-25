package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	DB  *gorm.DB
	app *fiber.App
}

func (s *Server) New(db *gorm.DB) *Server {
	app := fiber.New()

	s.DB = db
	s.app = app

	return s
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) ConnectRoutes() {
	s.app.Get("/", s.routeIndex)
}
