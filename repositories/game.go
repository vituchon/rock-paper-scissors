package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/vituchon/escobita/util"
)

type Game struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Owner        Player   `json:"owner"`
	Players      []Player `json:"players"`
	Started      bool     `json:"started"`
	CurrentMatch Match    `json:"currentMatch"`
}

type GameAction struct {
	Player Player `json:"player"`
	Weapon string `json:"weapon"`
}

type PlayerMove struct {
	Weapon string `json:"weapon"`
	When   int64  `json:"when"`
}

type Match struct {
	MoveByPlayerId map[int]PlayerMove `json:"lastMoveByPlayerId"`
}

func NewGame(players []Player) Game {
	return Game{
		CurrentMatch: Match{
			MoveByPlayerId: make(map[int]PlayerMove),
		},
	}
}

func (game Game) IsStarted() bool {
	return game.Started
}

func (game Game) IsMatchInProgress() bool {
	return len(game.CurrentMatch.MoveByPlayerId) > 0
}

func (game *Game) CreateNewMatch() {
	game.CurrentMatch = Match{
		MoveByPlayerId: make(map[int]PlayerMove),
	}
}

func (game Game) HasNoMovesInCurrentMatch() bool {
	return len(game.CurrentMatch.MoveByPlayerId) == 0
}

const MAX_PLAYER_PER_GAME int = 2

func (game Game) IsMatchCompleted() bool {
	return len(game.CurrentMatch.MoveByPlayerId) == MAX_PLAYER_PER_GAME
}

// PRE COND: match must be completed!
func (game Game) ResolveMatch() int {
	var player1Id, player2Id int
	var weapon1, weapon2 string
	for playerId, weapon := range game.CurrentMatch.MoveByPlayerId {
		if player1Id == 0 {
			player1Id = playerId
			weapon1 = weapon.Weapon
		} else {
			player2Id = playerId
			weapon2 = weapon.Weapon
		}
	}

	if weapon1 == weapon2 {
		return 0
	}
	if weapon1 == "✊" && weapon2 == "✌️" {
		return player1Id
	}
	if weapon1 == "✋" && weapon2 == "✊" {
		return player1Id
	}
	if weapon1 == "✌️" && weapon2 == "✋" {
		return player1Id
	}
	return player2Id
}

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

func (game *Game) StartGame() error {
	if game.IsStarted() {
		return GameStartedErr
	}
	if game.IsMatchInProgress() {
		return MatchInProgressErr
	}
	game.Started = true
	return nil
}

func (game *Game) PerformAction(action GameAction) error {
	if game.IsMatchCompleted() {
		err := fmt.Errorf("Can not perform  action: match in game(id='%d') is completed. %w", game.Id, MatchIsCompletedErr)
		return err
	}
	game.CurrentMatch.MoveByPlayerId[action.Player.Id] = PlayerMove{
		Weapon: action.Weapon,
		When:   time.Now().Unix(),
	}
	return nil
}

func (game Game) CanPlayerDeleteGame(player Player) bool {
	return game.Owner.Id == player.Id
}

var MatchInProgressErr error = errors.New("The match is in progress")
var MatchIsCompletedErr error = errors.New("The match is completed, can not take further actions.")
var GameStartedErr error = errors.New("The game is started")
var PlayerAlreadyJoinedErr error = errors.New("The player has already joined the game")
var PlayerNotJoinedErr error = errors.New("The player has not joined the game")
