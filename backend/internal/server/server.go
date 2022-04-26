package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	DB      *gorm.DB
	app     *fiber.App
	pending map[common.KickaeID][]*common.PendingPlayer
}

func New(db *gorm.DB) (s *Server) {
	app := fiber.New()
	s = &Server{
		DB:      db,
		app:     app,
		pending: make(map[common.KickaeID][]*common.PendingPlayer),
	}
	s.app.Get("/", s.routeIndex)
	s.app.Post("/e/rfid/:kicker_id/:goal_id/:player_id", s.routeRFID)
	s.app.Post("/e/tor/:kicker_id/:goal_id", s.routeTor)
	s.app.Get("/p/monitor/:kicker_id", s.routePull)
	return
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// db functions

func (s *Server) findActiveGameByKicker(id common.KickaeID) (g *common.Game, err error) {
	err = s.DB.Model(&common.Game{}).
		Where("kickae_id = ? AND end_time IS NULL", id).
		First(&g).Error
	return
}

func (s *Server) findPlayerById(id common.UserID) (p *common.Player, err error) {
	err = s.DB.Where(&common.Player{
		ID: id,
	}).First(&p).Error
	return
}
