package repositories

import (
	"errors"
)

var EntityNotExistsErr error = errors.New("Entity doesn't exists")
var DuplicatedEntityErr error = errors.New("Duplicated Entity")
var InvalidEntityStateErr error = errors.New("Entity state is invalid")

type Games interface {
	GetGames() ([]Game, error)
	GetGameById(id int) (*Game, error)
	CreateGame(game Game) (created *Game, err error)
	UpdateGame(game Game) (updated *Game, err error)
	DeleteGame(id int) error
	GetGamesCreatedCount(playerId int) int
}

type Players interface {
	GetPlayers() ([]Player, error)
	GetPlayerById(id int) (*Player, error)
	CreatePlayer(player Player) (created *Player, err error)
	UpdatePlayer(player Player) (updated *Player, err error)
	DeletePlayer(id int) error
}

type Messages interface {
	GetMessages() ([]PersistentMessage, error)
	GetMessagesByGame(gameId int) ([]PersistentMessage, error)
	GetMessagesByGameAndTime(gameId int, since int64) ([]PersistentMessage, error)
	GetMessageById(id int) (*PersistentMessage, error)
	CreateMessage(message PersistentMessage) (created *PersistentMessage, err error)
	UpdateMessage(message PersistentMessage) (updated *PersistentMessage, err error)
	DeleteMessage(id int) error
}
