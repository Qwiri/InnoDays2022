package internal

import "time"

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
	ID       uint
	Stime    time.Time
	Etime    *time.Time
	Sb       uint
	Sw       uint
	KickaeID KickaeID

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
	Games []Game
}

type Player struct {
	ID  string
	elo uint

	Games []Game `gorm:"many2many:GamePlayers"`
}
