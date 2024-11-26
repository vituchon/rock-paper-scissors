package repositories

import (
	"errors"
	"sync"

	"github.com/vituchon/escobita/util"
)

type Game struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Owner   Player   `json:"owner"`
	Players []Player `json:"players"`
	Started bool     `json:"started"`
}

func NewGame(players []Player) Game {
	return Game{}
}

func (game Game) IsStarted() bool {
	return game.Started
}

var MatchInProgressErr error = errors.New("The match is in progress")
var GameStartedErr error = errors.New("The game is started")
var PlayerAlreadyJoinedErr error = errors.New("The player has already joined the game")
var PlayerNotJoinedErr error = errors.New("The player has not joined the game")

func (game *Game) Join(player Player) error {
	if game.IsStarted() {
		return GameStartedErr
	}
	joinedPlayer := util.Find(game.Players, func(gamePlayer Player) bool { return gamePlayer.Id == player.Id })
	playerNotJoined := joinedPlayer == nil
	if playerNotJoined {
		game.Players = append(game.Players, player)
	} else {
		return PlayerAlreadyJoinedErr
	}
	return nil
}

func (game *Game) IsJoined(player Player) bool {
	joinedPlayer := util.Find(game.Players, func(gamePlayer Player) bool { return gamePlayer.Id == player.Id })
	return joinedPlayer != nil
}

func (game *Game) Quit(player Player) error {
	if game.IsStarted() {
		return GameStartedErr // if the game is already started you can quit! as The Eagles stated "you can checkout any time you want, but you can never leave"
	}
	var playerIndex int = -1
	for i, gamePlayer := range game.Players {
		if gamePlayer.Id == player.Id {
			playerIndex = i
			break
		}
	}
	playerJoined := playerIndex != -1
	if playerJoined {
		game.Players = append(game.Players[:playerIndex], game.Players[playerIndex+1:]...) // taken advice from https://github.com/golang/go/wiki/SliceTricks#delete
	} else {
		return PlayerNotJoinedErr
	}
	return nil
}

type GamesMemoryRepository struct {
	gamesById              map[int]Game
	gamesCreatedByPlayerId map[int]int
	idSequence             int
	mutex                  sync.Mutex
}

func NewGamesMemoryRepository() *GamesMemoryRepository {
	return &GamesMemoryRepository{gamesById: make(map[int]Game), gamesCreatedByPlayerId: make(map[int]int), idSequence: 0}
}

func (repo *GamesMemoryRepository) GetGames() ([]Game, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	games := make([]Game, 0, len(repo.gamesById))
	for _, game := range repo.gamesById {
		games = append(games, game)
	}
	return games, nil
}

func (repo *GamesMemoryRepository) GetGameById(id int) (*Game, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	game, exists := repo.gamesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &game, nil
}

func (repo *GamesMemoryRepository) CreateGame(game Game) (created *Game, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	nextId := repo.idSequence + 1 // need to copy repo.idSequence in another place (nextId), also added plus one to increment the sequence current number...
	game.Id = nextId              // ...that place (nextId) will work as reference...
	repo.gamesById[nextId] = game
	repo.idSequence++ // ...if it is used idSequence as a reference, then each update would increment all the games Id by 1 (actually all game.Id will point to the same thing)
	repo.gamesCreatedByPlayerId[game.Owner.Id]++
	return &game, nil
}

func (repo *GamesMemoryRepository) UpdateGame(game Game) (updated *Game, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.gamesById[game.Id] = game
	return &game, nil
}

func (repo *GamesMemoryRepository) DeleteGame(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	game := repo.gamesById[id]
	repo.gamesCreatedByPlayerId[game.Owner.Id]--
	delete(repo.gamesById, id)
	return nil
}

func (repo GamesMemoryRepository) GetGamesCreatedCount(playerId int) int {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	return repo.gamesCreatedByPlayerId[playerId]
}
