package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func (s *Server) routeTor(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID internal.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse kicker_id")
	} else {
		kickerID = internal.KickaeID(k)
	}

	// parse goal id
	var goalID internal.TeamColor
	if g, err := c.ParamsInt("goal_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse goal_id")
	} else {
		goalID = internal.TeamColor(g)
	}

	// check if the kicker has a currently running game
	g := &internal.Game{}
	if err = s.DB.Where(&internal.Game{
		Etime:    nil,
		KickaeID: kickerID,
	}).First(g).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	// if no game is currently running
	if err == gorm.ErrRecordNotFound {
		if len(s.pendingPlayers[kickerID]) < 2 {
			return fiber.NewError(fiber.StatusConflict, "Can't start game with less than two players")
		}
		// TODO: check if logged in same team
		// TODO: check if same player

		g = &internal.Game{
			Stime:    time.Now(),
			KickaeID: kickerID,
		}
		if err = s.DB.Create(g).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		gps := s.pendingPlayers[kickerID]
		for _, gp := range gps {
			gp.GameID = g.ID

			if err = s.DB.Create(gp).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}
		delete(s.pendingPlayers, kickerID)

	} else {
		if goalID == internal.BlackTeamColor {
			g.Sw++
		} else if goalID == internal.WhiteTeamColor {
			g.Sb++
		}
		if err = s.DB.Where("ID = ?", g.ID).Updates(g).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

	}

	return c.SendStatus(fiber.StatusCreated)
}
