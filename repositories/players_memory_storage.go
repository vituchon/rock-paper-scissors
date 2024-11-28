package repositories

import (
	"sync"
)

type Player struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Emotar string `json:"emotar"`
}

type PlayersMemoryRepository struct {
	playersById map[int]Player
	mutex       sync.Mutex
}

func NewPlayersMemoryRepository() *PlayersMemoryRepository {
	return &PlayersMemoryRepository{playersById: make(map[int]Player)}
}

func (repo *PlayersMemoryRepository) GetPlayers() ([]Player, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	players := make([]Player, 0, len(repo.playersById))
	for _, player := range repo.playersById {
		players = append(players, player)
	}
	return players, nil
}

func (repo *PlayersMemoryRepository) GetPlayerById(id int) (*Player, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	player, exists := repo.playersById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &player, nil
}

func (repo *PlayersMemoryRepository) CreatePlayer(player Player) (created *Player, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.playersById[player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryRepository) UpdatePlayer(player Player) (updated *Player, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.playersById[player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryRepository) DeletePlayer(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	delete(repo.playersById, id)
	return nil
}
