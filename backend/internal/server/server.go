package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	DB             *gorm.DB
	app            *fiber.App
	pendingPlayers map[common.KickaeID][]*common.GamePlayers
}

func New(db *gorm.DB) (s *Server) {
	app := fiber.New()
	s = &Server{
		DB:             db,
		app:            app,
		pendingPlayers: make(map[common.KickaeID][]*common.GamePlayers),
	}
	s.app.Get("/", s.routeIndex)
	s.app.Post("/e/rfid/:kicker_id/:goal_id/:player_id", s.routeRFID)
	s.app.Post("/e/tor/:kicker_id/:goal_id", s.routeTor)
	return
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
