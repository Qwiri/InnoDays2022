package server

import (
	"errors"
	"fmt"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func (s *Server) routeRFID(c *fiber.Ctx) (err error) {
	// parse kicker id
	var kickerID common.KickaeID
	if k, err := c.ParamsInt("kicker_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing kicker_id failed")
	} else {
		kickerID = common.KickaeID(k)
	}
	// check if there's already a game running on kickerID
	if _, err := s.findActiveGameByKicker(kickerID); err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotAcceptable, "game already running")
	}

	// parse goal id
	var goalID common.GoalColor
	if g, err := c.ParamsInt("goal_id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parsing goal_id failed")
	} else {
		goalID = common.GoalColor(g)
	}
	// check if goalID is valid
	if goalID != common.WhiteTeamColor && goalID != common.BlackTeamColor {
		return fiber.NewError(fiber.StatusBadRequest, "invalid goal ID")
	}

	// parse player id
	playerID := common.UserID(utils.CopyString(c.Params("player_id")))
	// even through the playerID is a string, it should only contain numbers
	if !common.UserIDPattern.MatchString(string(playerID)) {
		return fiber.NewError(fiber.StatusBadRequest, "user id invalid")
	}

	// check if player is already pending on another team
	// if the player is already pending on another team, remove the player from that team and add to current team
	// if there are too many players in the current team, remove the oldest player from the team
	// WARNING: looking at the code below you might want to take a look at /r/eyebleach

	// check if there are any pending players on the current kicker table
	if pending := s.pending[kickerID]; len(pending) > 0 {
		var (
			co uint
			pe []*common.PendingPlayer
		)
		// remove player from any other teams
		for _, p := range pending {
			if p.PlayerID != playerID {
				pe = append(pe, p)
				// count the number of players on the current team
				if p.Team == goalID {
					co++
				}
			}
		}
		s.pending[kickerID] = pe

		// check if there are too many players on the current team
		// if there are already 2 players (or more) on the current team,
		// remove the oldest player from the team
		if co >= 2 {
			var (
				oldest       time.Duration
				oldestPlayer *common.PendingPlayer
			)
			// find oldest player
			for _, p := range pending {
				if p.Team != goalID {
					continue
				}
				if oldestPlayer == nil {
					oldestPlayer = p
					oldest = p.AddedAt.Sub(time.Now())
				} else {
					if s := p.AddedAt.Sub(time.Now()); s > oldest {
						oldestPlayer = p
						oldest = s
					}
				}
			}
			// remove oldestPlayer from pending players
			if oldestPlayer != nil {
				pe = nil
				for _, p := range pending {
					if p.PlayerID != oldestPlayer.PlayerID {
						pe = append(pe, p)
					}
				}
				s.pending[kickerID] = pe
			}
		}
	}

	// create pending player
	pending := &common.PendingPlayer{
		PlayerID: playerID,
		Team:     goalID,
		AddedAt:  time.Now(),
	}

	// add to pending players
	s.pending[kickerID] = append(s.pending[kickerID], pending)
	fmt.Println(s.pending[kickerID])

	return c.Status(fiber.StatusCreated).
		SendString(string(playerID) + " -> team " + strconv.Itoa(int(goalID)))
}
