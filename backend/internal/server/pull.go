package server

import (
	"errors"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PullPayload struct {
	Game    *common.Game
	Pending []struct {
		Player  *common.Player
		Pending *common.PendingPlayer
	}
}

func (s *Server) routePull(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse kicker_id")
	} else {
		kickerID = common.KickaeID(k)
	}

	var (
		p    PullPayload
		game common.Game
	)

	// check if the kicker has a currently running game
	if err = s.DB.
		Preload("Players").
		Preload("Players.Player").
		Preload("Goals").
		Where("kickae_id = ? AND end_time IS NULL", kickerID).
		First(&game).Error; err != nil {

		// ErrRecordNotFound is okay. A game hasn't been started yet.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c = c.Status(fiber.StatusNotFound)
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		p.Game = &game
		c = c.Status(fiber.StatusAccepted)
	}

	// pending
	if pending, ok := s.pending[kickerID]; ok {
		for _, pen := range pending {
			player := s.getPlayerById(pen.PlayerID)
			p.Pending = append(p.Pending, struct {
				Player  *common.Player
				Pending *common.PendingPlayer
			}{Player: &player, Pending: pen})
		}
	}

	return c.JSON(&p)
}
