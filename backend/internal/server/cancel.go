package server

import (
	"errors"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (s *Server) routeCancel(c *fiber.Ctx) (err error) {
	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse kicker_id")
	} else {
		kickerID = common.KickaeID(k)
	}

	// check if the kicker has a currently running game
	var game *common.Game
	if game, err = s.findActiveGameByKicker(kickerID); err != nil {
		// unknown error - return
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusNotFound, "no active game found")
	}

	if err = game.End(s.DB, common.ReasonCancel); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("game canceled")
}
