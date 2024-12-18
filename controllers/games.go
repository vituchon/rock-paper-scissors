package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/vituchon/rock-paper-scissors/repositories"
	"github.com/vituchon/rock-paper-scissors/services"
)

var gamesRepository repositories.Games = repositories.NewGamesMemoryRepository()

func GetGames(response http.ResponseWriter, request *http.Request) {
	games, err := gamesRepository.GetGames()
	if err != nil {
		msg := fmt.Sprintf("error while retrieving games : '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, games)
}

func GetGameById(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game(id=%d): '%v'", id, err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, game)
}

const MAX_GAMES_PER_PLAYER = 1

func CreateGame(response http.ResponseWriter, request *http.Request) {
	playerId := services.GetClientId(request) // will be the game's owner
	if gamesRepository.GetGamesCreatedCount(playerId) == MAX_GAMES_PER_PLAYER {
		msg := fmt.Sprintf("Player(id='%d') has reached the maximum game creation limit: '%v'", playerId, MAX_GAMES_PER_PLAYER)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByValue(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		msg := fmt.Sprintf("error getting player by id='%d': '%v'", playerId, err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game.Owner = *player
	created, err := gamesRepository.CreateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while creating game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdateGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByValue(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while updating game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingActionMsgPayload{updated, nil}
	services.GameWebSockets.NotifyGameConns(game.Id, "updated", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeleteGame(response http.ResponseWriter, request *http.Request) {
	var player repositories.Player
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		msg := fmt.Sprintf("error reading request body: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	id := game.Id
	err = gamesRepository.DeleteGame(id)
	if err != nil {
		msg := fmt.Sprintf("error while deleting game(id=%d): '%v'", id, err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	services.GameWebSockets.UnbindAllWebSocketsInGame(id, request)
	response.WriteHeader(http.StatusOK)
}

func DeleteGames(response http.ResponseWriter, request *http.Request) {
	games, err := gamesRepository.GetGames()
	if err != nil {
		msg := fmt.Sprintf("error while retrieving games : '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	for _, game := range games {
		gamesRepository.DeleteGame(game.Id)
		if err != nil {
			msg := fmt.Sprintf("error while deleting game(id=%d): '%v'", game.Id, err)
			log.Println(msg)
			http.Error(response, msg, http.StatusInternalServerError)
			return
		}
		services.GameWebSockets.UnbindAllWebSocketsInGame(game.Id, request)
	}
	response.WriteHeader(http.StatusOK)
}

func StartGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	playerId := services.GetClientId(request)
	if game.Owner.Id != playerId {
		msg := fmt.Sprintf("error while starting game: request doesn't cames from the owner, in cames from %d", playerId)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	err = game.StartGame()
	if err != nil {
		msg := fmt.Sprintf("error while starting game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	game.CreateNewMatch()

	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while starting game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	msgPayload := services.WebSockectOutgoingActionMsgPayload{game, nil}
	services.GameWebSockets.NotifyGameConns(game.Id, "game-start", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func RestartGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	playerId := services.GetClientId(request)
	if game.Owner.Id != playerId {
		msg := fmt.Sprintf("error while restarting game: request doesn't cames from the owner, in cames from %d", playerId)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	if game.HasNoMovesInCurrentMatch() {
		msg := fmt.Sprintf("error while restarting game: match has no moves. Skipping restart.")
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game.CreateNewMatch()
	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while restarting game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	msgPayload := services.WebSockectOutgoingActionMsgPayload{game, nil}
	services.GameWebSockets.NotifyGameConns(game.Id, "game-restart", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func JoinGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	playerId := services.GetClientId(request)
	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		msg := fmt.Sprintf("error while getting player by id, error was: '%v'", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	err = game.Join(*player)
	if err != nil {
		msg := fmt.Sprintf("error while joining game, error was: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while updating game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingJoinMsgPayload{updated, player}
	services.GameWebSockets.NotifyGameConns(game.Id, "player-join", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func PerformAction(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return

	}
	var action repositories.GameAction
	err = parseJsonFromReader(request.Body, &action)
	if err != nil {
		msg := fmt.Sprintf("error reading request body: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	err = game.PerformAction(action)
	if err != nil {
		msg := fmt.Sprintf("error while performing action: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while updating game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingActionMsgPayload{game, &action}
	services.GameWebSockets.NotifyGameConns(game.Id, "player-action", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func ResolveCurrentGameMatch(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	playerId := services.GetClientId(request)
	if game.Owner.Id != playerId {
		msg := fmt.Sprintf("error while resolving current game's match: request doesn't cames from the owner, in cames from %d", playerId)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	winnerPlayerId := game.ResolveMatch()
	/*game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while updating game: '%v'", err)
		http.Error(response,msg,http.StatusInternalServerError)
		return
	}*/
	msgPayload := &winnerPlayerId
	services.GameWebSockets.NotifyGameConns(game.Id, "game-match-resolved", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func QuitGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return
	}

	playerId := services.GetClientId(request)
	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		msg := fmt.Sprintf("error while getting player by id, error was: '%v'", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	err = game.Quit(*player)
	if err != nil {
		msg := fmt.Sprintf("error while quiting game, error was: '%v'", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		msg := fmt.Sprintf("error while updating game: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingJoinMsgPayload{updated, player}
	services.GameWebSockets.NotifyGameConns(game.Id, "game-quit", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

type WebSockectOutgoingChatMsgPayload struct {
	Message services.VolatileWebMessage `json:"message"`
}

func SendMessage(response http.ResponseWriter, request *http.Request) {
	var message services.VolatileWebMessage
	err := parseJsonFromReader(request.Body, &message)
	if err != nil {
		msg := fmt.Sprintf("error reading request body: '%v'", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByReference(request)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving game: '%v'", err)
		log.Println(msg)
		if errors.Is(err, NoGameIdRouteParamErr) {
			http.Error(response, msg, http.StatusBadRequest)
		} else {
			http.Error(response, msg, http.StatusInternalServerError)
		}
		return

	}
	id := game.Id
	msgPayload := WebSockectOutgoingChatMsgPayload{message}
	services.GameWebSockets.NotifyGameConns(id, "game-chat", msgPayload)

	WriteJsonResponse(response, http.StatusOK, struct{}{})
}

func BindClientWebSocketToGame(response http.ResponseWriter, request *http.Request) {
	gameId, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	services.GameWebSockets.BindClientWebSocketToGame(response, request, gameId)
	response.WriteHeader(http.StatusOK)
}

func UnbindClientWebSocketInGame(response http.ResponseWriter, request *http.Request) {
	conn := services.WebSocketsHandler.Retrieve(request)
	if conn != nil {
		services.GameWebSockets.UnbindClientWebSocketInGame(conn, request)
		response.WriteHeader(http.StatusOK)
	} else {
		msg := fmt.Sprintf("No need to release web socket as it was not adquired (or already released) for  client(id='%d')", services.GetClientId(request))
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
	}
}

var NoGameIdRouteParamErr = errors.New("the request URL is missing the game ID as a route parameter")

// Retrieves the stored game in the underlying storage system using the id present in the URL (route param)
func retrieveGameByReference(request *http.Request) (*repositories.Game, error) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		err = fmt.Errorf("%w: %v", NoGameIdRouteParamErr, err)
		return nil, err
	}

	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		errMsg := fmt.Sprintf("error while retrieving game(id=%d): '%v'", id, err)
		return nil, errors.New(errMsg)
	}
	return game, nil
}

// Retrieves the stored game in the request's payload
func retrieveGameByValue(request *http.Request) (*repositories.Game, error) {
	var game repositories.Game
	/*bufferedReader := bufio.NewReader(request.Body)

	// Read all data into a single buffer
	buffer, err := bufferedReader.ReadBytes(0) // 0 means to read until the end
	if err != nil && err != io.EOF {
		errMsg := fmt.Sprintf("Error reading from reader: %v", er)
		return nil, errors.New(errMsg)
	}

	// Print the entire content
	fmt.Println("Game:", string(buffer))

	err = parseJsonFromReader(bytes.NewReader(buffer), &game)*/

	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		errMsg := fmt.Sprintf("error reading request body: '%v'", err)
		return nil, errors.New(errMsg)
	}
	return &game, nil
}
