package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	DB             *gorm.DB
	app            *fiber.App
	pendingPlayers map[internal.KickaeID][]*internal.GamePlayers
}

func (s *Server) New(db *gorm.DB) *Server {
	app := fiber.New()

	s.DB = db
	s.app = app
	s.pendingPlayers = make(map[internal.KickaeID][]*internal.GamePlayers)

	return s
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) ConnectRoutes() {
	s.app.Get("/", s.routeIndex)
	s.app.Post("/e/rfid/:kicker_id/:goal_id/:player_id", s.routeRFID)
	s.app.Post("/e/tor/:kicker_id/:goal_id", s.routeTor)
}
