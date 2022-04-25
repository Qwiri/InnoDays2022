package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) routeRFID(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID internal.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing kicker_id failed")
	} else {
		kickerID = internal.KickaeID(k)
	}

	// parse goal id
	goalID, err := c.ParamsInt("goal_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing goal_id failed")
	}

	// parse player id
	playerID := c.Params("player_id")

	p := &internal.GamePlayers{
		PlayerID: playerID,
		Team:     internal.TeamColor(goalID),
	}

	// add to kickae
	s.pendingPlayers[kickerID] = append(s.pendingPlayers[kickerID], p)

	return c.SendStatus(fiber.StatusCreated)
}
