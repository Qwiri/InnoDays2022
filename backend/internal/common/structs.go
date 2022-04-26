package common

import (
	"database/sql"
	"regexp"
	"time"
)

type UserID string

var UserIDPattern = regexp.MustCompile("^\\d+$")

type GoalColor uint
type KickaeID uint

const (
	BlackTeamColor GoalColor = iota + 1
	WhiteTeamColor
)

var TableModels = []interface{}{
	&Game{},
	&Kickae{},
	&Player{},
}

type Game struct {
	ID         uint
	StartTime  time.Time
	EndTime    sql.NullTime
	ScoreBlack uint
	ScoreWhite uint
	KickaeID   KickaeID

	Players []*Player `gorm:"many2many:game_players"`
}

type PendingPlayer struct {
	PlayerID UserID
	Team     GoalColor
	AddedAt  time.Time // used for janitor
}

type Kickae struct {
	ID    uint
	Room  string
	Note  string
	Games []*Game
}

type Player struct {
	ID  UserID
	Elo uint

	Games []Game `gorm:"many2many:GamePlayers"`
}
