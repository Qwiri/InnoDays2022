package common

import (
	"database/sql"
	"regexp"
	"time"
)

type (
	UserID    string
	GoalColor uint
	KickaeID  uint
)

const (
	BlackTeamColor GoalColor = iota + 1
	WhiteTeamColor
)

var UserIDPattern = regexp.MustCompile("^\\d+$")

// TableModels contains all the models that are used in the database.
// used for GORM AutoMigrate
var TableModels = []interface{}{
	&Game{},
	&Kickae{},
	&Player{},
	&GamePlayers{},
}

// Game represents a game
type Game struct {
	// ID of the game
	ID uint
	// StartTime is the time when the game started
	StartTime time.Time
	// EndTime is the time when the game ended - null by default
	EndTime sql.NullTime
	// ScoreBlack is the score of the black team
	ScoreBlack uint
	// ScoreWhite is the score of the white team
	ScoreWhite uint
	// KickaeID is the id of the kickae the game is taking place
	KickaeID KickaeID
	// Kickae is the kickae, where the game is taking place at
	Kickae Kickae `gorm:"foreignKey:KickaeID"`
	// UpdatedAt is the time when the game was last updated - used for janitor
	UpdatedAt time.Time
	// Players are the players in the game
	Players []*GamePlayers
}

// GamePlayers represents a player in a game
type GamePlayers struct {
	// PlayerID is the id of the player
	PlayerID UserID `gorm:"primaryKey"`
	// Player is the player
	Player *Player `gorm:"foreignKey:PlayerID"`

	// GameID is the id of the game
	GameID uint `gorm:"primaryKey"`
	// Game is the game
	Game *Game `gorm:"foreignKey:GameID"`

	// Team is the goal color of the player (black/white)
	Team GoalColor
}

// Kickae represents a kicker
type Kickae struct {
	// ID of the kicker
	ID uint
	// Room is the room the kicker is in
	Room string
	// Note is an optional note of the kicker
	Note string
	// Games is the games the kicker is in
	Games []*Game
}

// Player represents a player
type Player struct {
	// ID of the player
	ID UserID
	// Nick is the nickname of the player - optional; can be empty
	Nick string
	// Elo is the elo of the player
	Elo uint
}

// PendingPlayer represents a player that is waiting to join a game
type PendingPlayer struct {
	// PlayerID is the id of the player
	PlayerID UserID
	// Team is the goal color of the player (black/white)
	Team GoalColor
	// AddedAt is the time when the player was added - used for janitor
	AddedAt time.Time
}
