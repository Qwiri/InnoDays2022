package server

import (
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func (s *Server) routeTor(c *fiber.Ctx) (err error) {

	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse kicker_id")
	} else {
		kickerID = common.KickaeID(k)
	}

	// parse goal id
	var goalID common.TeamColor
	if g, err := c.ParamsInt("goal_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse goal_id")
	} else {
		goalID = common.TeamColor(g)
	}

	// check if the kicker has a currently running game
	g := &common.Game{}
	if err = s.DB.Where(&common.Game{
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

		g = &common.Game{
			StartTime: time.Now(),
			KickaeID:  kickerID,
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
		if goalID == common.BlackTeamColor {
			g.ScoreWhite++
		} else if goalID == common.WhiteTeamColor {
			g.ScoreBlack++
		}
		if err = s.DB.Where("ID = ?", g.ID).Updates(g).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

	}

	return c.SendStatus(fiber.StatusCreated)
}
