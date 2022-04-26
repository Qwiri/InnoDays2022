package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PullPayload struct {
	Game *common.Game
}

func (s *Server) routePull(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse kicker_id")
	} else {
		kickerID = common.KickaeID(k)
	}

	p := PullPayload{}

	// check if the kicker has a currently running game
	var game common.Game
	if err = s.DB.
		Preload("Players").
		Preload("Players.Player").
		Where(&common.Game{
			KickaeID: kickerID,
		}).
		First(&game).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	p.Game = &game
	return c.JSON(&p)
}
