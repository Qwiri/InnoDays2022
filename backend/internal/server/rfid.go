package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) routeRFID(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing kicker_id failed")
	} else {
		kickerID = common.KickaeID(k)
	}

	// parse goal id
	goalID, err := c.ParamsInt("goal_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing goal_id failed")
	}

	// parse player id
	playerID := c.Params("player_id")

	p := &common.GamePlayers{
		PlayerID: playerID,
		Team:     common.TeamColor(goalID),
	}

	// add to kickae
	s.pendingPlayers[kickerID] = append(s.pendingPlayers[kickerID], p)

	return c.SendStatus(fiber.StatusCreated)
}
