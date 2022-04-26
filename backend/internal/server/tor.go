package server

import (
	"errors"
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
	var goalID common.GoalColor
	if g, err := c.ParamsInt("goal_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not parse goal_id")
	} else {
		goalID = common.GoalColor(g)
	}

	// check if the kicker has a currently running game
	var game *common.Game
	if game, err = s.findActiveGameByKicker(kickerID); err != nil {
		// unknown error - return
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// check if there are enough pending players on both teams
		pt := make(map[common.GoalColor]int)
		for _, p := range s.pending[kickerID] {
			pt[p.Team]++
		}
		if pt[common.BlackTeamColor] <= 0 || pt[common.WhiteTeamColor] <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "not enough players on both teams")
		}

		// collect players from pending list and load them from pending players table
		var players []*common.GamePlayers
		for _, p := range s.pending[kickerID] {
			// find player
			var player *common.Player
			if player, err = s.findPlayerById(p.PlayerID); err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
				// create and save player
				player = &common.Player{
					ID: p.PlayerID,
				}
				if err = s.DB.Create(player).Error; err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}
			players = append(players, &common.GamePlayers{
				PlayerID: player.ID,
				GameID:   game.ID,
				Team:     p.Team,
			})
		}

		// create new game object
		game = &common.Game{
			StartTime: time.Now(),
			KickaeID:  kickerID,
			Players:   players,
		}
		if err = s.DB.Create(game).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// remove pending players
		delete(s.pending, kickerID)
	} else {
		// game found
		if goalID == common.BlackTeamColor {
			game.ScoreWhite++
		} else if goalID == common.WhiteTeamColor {
			game.ScoreBlack++
		}
		if err = s.DB.Where("ID = ?", game.ID).Updates(game).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}
