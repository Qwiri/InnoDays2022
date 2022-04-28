package server

import (
	"errors"
	"fmt"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
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
	if game, err = s.findActiveGameByKicker(kickerID, true); err != nil {
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
			player := s.getPlayerById(p.PlayerID)
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

		return c.Status(fiber.StatusCreated).SendString("game created")
	} else {
		// game found
		if goalID == common.BlackTeamColor {
			game.ScoreWhite++
		} else if goalID == common.WhiteTeamColor {
			game.ScoreBlack++
		}

		// add goal
		if err = s.DB.Create(&common.Goal{
			GameID: game.ID,
			Team:   goalID,
			Time:   time.Now(),
		}).Error; err != nil {
			log.WithError(err).Warn("cannot save goal")
		}

		// add score to game
		if err = s.DB.Where("ID = ?", game.ID).Updates(game).Error; err != nil {
			log.WithError(err).Warn("cannot update game")
		}

		var (
			winner common.GoalColor
			ww     = won(game.ScoreWhite, game.ScoreBlack)
			wb     = won(game.ScoreBlack, game.ScoreWhite)
		)

		// check game
		if ww || wb {
			if ww {
				winner = common.WhiteTeamColor
			} else if wb {
				winner = common.BlackTeamColor
			} else {
				return fiber.NewError(fiber.StatusInternalServerError, "cannot determine winner")
			}
			log.WithField("winner", winner).Info("game ended")

			// game is over
			if err = game.End(s.DB, common.ReasonWin); err != nil {
				log.WithError(err).Warn("cannot end game")
			}

			// calculate points and add them to players
			s.elo(*game, winner)

			return c.Status(fiber.StatusCreated).SendString("game won")
		}

		// send ok
		return c.Status(fiber.StatusAccepted).SendString("count goal")
	}
}

func won(a, b uint) bool {
	return a >= 10 && math.Abs(float64(a)-float64(b)) >= 2
}

func (s *Server) elo(game common.Game, winner common.GoalColor) {
	const (
		u = 400.0
		k = 40
	)

	// normalize scores
	if game.ScoreBlack > 10 && game.ScoreBlack > game.ScoreWhite {
		game.ScoreBlack = 10
		game.ScoreWhite = 8
	} else if game.ScoreWhite > 10 && game.ScoreWhite > game.ScoreBlack {
		game.ScoreWhite = 10
		game.ScoreBlack = 8
	}

	var (
		sB float64
		sW float64
	)
	if winner == common.BlackTeamColor {
		sB = 10.0 / float64(10+game.ScoreWhite)
		sW = 1 - sB
	} else {
		sW = 10.0 / float64(10+game.ScoreBlack)
		sB = 1 - sW
	}

	var (
		mEW float64
		mEB float64
	)
	{
		elos := make(map[common.GoalColor][]uint)
		for _, p := range game.Players {
			elos[p.Team] = append(elos[p.Team], p.Player.Elo)
		}
		mw := make(map[common.GoalColor]float64)
		for t, e := range elos {
			var sum uint
			for _, s := range e {
				sum += s
			}
			mw[t] = float64(sum) / float64(len(e))
		}
		mEW = mw[common.WhiteTeamColor]
		mEB = mw[common.BlackTeamColor]

		fmt.Println("elos:", elos)
	}
	fmt.Println("mEW:", mEW, "mEB:", mEB)

	var (
		eW = 1 / (1 + math.Pow(10, (mEB-mEW)/u))
		eB = 1 / (1 + math.Pow(10, (mEW-mEB)/u))
	)
	fmt.Println("eW:", eW, "eB:", eB)

	for _, p := range game.Players {
		var (
			syu float64
			e   float64
		)
		if p.Team == common.WhiteTeamColor {
			syu = sW
			e = eW
		} else {
			syu = sB
			e = eB
		}
		nE := int(math.Ceil(float64(p.Player.Elo) + k*(syu-e)))
		if nE < 0 {
			nE = 0
		}

		// set new elo
		el := &common.EloLog{
			PlayerID: p.PlayerID,
			OldElo:   p.Player.Elo,
			NewElo:   uint(nE),
			Time:     time.Now(),
		}

		// update elo
		if err := s.DB.Model(p.Player).
			Where(&common.Player{ID: p.PlayerID}).
			Updates(&common.Player{Elo: uint(nE)}).Error; err != nil {
			log.WithError(err).Warn("cannot update player's elo")
		}

		// insert log entry
		if err := s.DB.Create(el).Error; err != nil {
			log.WithError(err).Warn("cannot save log entry")
		}
	}
}
