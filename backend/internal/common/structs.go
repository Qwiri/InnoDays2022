package common

import (
	"database/sql"
	"time"
)

type UserID uint
type TeamColor uint
type KickaeID uint

const (
	BlackTeamColor TeamColor = iota + 1
	WhiteTeamColor
)

var TableModels = []interface{}{
	&Game{},
	&Kickae{},
	&Player{},
	&GamePlayers{},
}

type Game struct {
	ID         uint
	StartTime  time.Time
	EndTime    sql.NullTime
	ScoreBlack uint
	ScoreWhite uint
	KickaeID   KickaeID

	Players []Player `gorm:"many2many:game_players"`
}

type GamePlayers struct {
	GameID   uint
	PlayerID string
	Team     TeamColor
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
