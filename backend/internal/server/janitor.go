package server

import (
	"context"
	"github.com/Qwiri/InnoDays2022/backend/internal/common"
	"github.com/apex/log"
	"gorm.io/gorm"
	"time"
)

type Janitor struct {
	db     *gorm.DB
	delay  time.Duration
	ticker *time.Ticker
	cancel context.Context
	svr    *Server
}

func NewJanitor(db *gorm.DB, svr *Server, cancel context.Context, delay time.Duration) *Janitor {
	return &Janitor{
		db:     db,
		svr:    svr,
		delay:  delay,
		cancel: cancel,
	}
}

func (j *Janitor) Start() {
	log.Info("[Janitor] Starting")
	if j.ticker != nil {
		j.Stop()
	}
	j.ticker = time.NewTicker(j.delay)
	for {
		select {
		case <-j.cancel.Done():
			log.Info("Stopped Janitor")
			return
		case <-j.ticker.C:
			j.Clean()
		}
	}
}

func (j *Janitor) Clean() {
	log.Info("[Janitor] Clean")

	// pending players
	for kickaeID, pendingPlayers := range j.svr.pending {
		var pending []*common.PendingPlayer
		for _, pendingPlayer := range pendingPlayers {
			if time.Since(pendingPlayer.AddedAt) < time.Minute*5 {
				pending = append(pending, pendingPlayer)
			} else {
				log.WithFields(log.Fields{
					"kickaeID": kickaeID,
					"player":   pendingPlayer.PlayerID,
				}).Info("[Janitor] Removing pending player")
			}
		}
		if len(pendingPlayers) <= 0 {
			delete(j.svr.pending, kickaeID)
		} else {
			j.svr.pending[kickaeID] = pending
		}
	}

	// pending game
	var games []*common.Game
	if err := j.db.Where("end_time IS NULL").Find(&games).Error; err != nil {
		log.WithError(err).Error("[Janitor] Failed to load games")
		return
	}

	log.Infof("[Janitor] Found %d games", len(games))
	for _, game := range games {
		if time.Since(game.UpdatedAt) < time.Minute*10 {
			continue
		}
		log.WithFields(log.Fields{
			"gameID": game.ID,
		}).Info("[Janitor] Marking game as done")
		if err := j.db.Model(game).Where(game).Update("end_time", time.Now()).Error; err != nil {
			log.WithError(err).Error("[Janitor] Failed to delete game")
		}
	}
}

func (j *Janitor) Stop() {
	if j.ticker != nil {
		log.Info("[Janitor] Stop")
		j.ticker.Stop()
		j.ticker = nil
	}
}
