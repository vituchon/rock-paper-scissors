package repositories

import (
	"sync"
)

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
	nextId := repo.idSequence + 1
	game.Id = nextId
	repo.gamesById[nextId] = game
	repo.idSequence++
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
