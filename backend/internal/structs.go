package internal

import "time"

var TableModels = []interface{}{
	&Game{},
	&Kickae{},
	&Player{},
}

type Game struct {
	ID       uint
	Stime    time.Time
	Etime    *time.Time
	Sb       uint
	Sw       uint
	KickaeID uint
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
}
