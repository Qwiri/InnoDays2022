package server

import (
	"database/sql"
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
	return
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// db functions

func (s *Server) findGameByKicker(id common.KickaeID) (g *common.Game, err error) {
	err = s.DB.Where(&common.Game{
		KickaeID: id,
		EndTime:  sql.NullTime{Valid: true}, // where EndTime = NULL
	}).First(&g).Error
	return
}

func (s *Server) findPlayerById(id common.UserID) (p *common.Player, err error) {
	err = s.DB.Where(&common.Player{
		ID: id,
	}).First(&p).Error
	return
}
