package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/vituchon/rock-paper-scissors/repositories"
	"github.com/vituchon/rock-paper-scissors/services"
)

var playersRepository repositories.Players = repositories.NewPlayersMemoryRepository()

func GetPlayers(response http.ResponseWriter, request *http.Request) {
	players, err := playersRepository.GetPlayers()
	if err != nil {
		log.Printf("error while retrieving players : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, players)
}

var getClientPlayerMutex sync.Mutex

func RegisterPlayer(response http.ResponseWriter, request *http.Request) {
	getClientPlayerMutex.Lock()
	defer getClientPlayerMutex.Unlock()
	id := services.GetClientId(request)
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		if err == repositories.EntityNotExistsErr {
			name, err := ParseSingleStringUrlQueryParam(request, "name")
			emotar, err := ParseSingleStringUrlQueryParam(request, "emotar")
			player, err = createPlayer(id, *name, *emotar)
			if err != nil {
				log.Printf("error while registering (create) client player : '%v'", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = playersRepository.CreatePlayer(*player) // saves player in a persistent storage
			if err != nil {
				log.Printf("error while registering (create) client player : '%v'", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("Creating new player %+v for ip=%s\n", player, request.RemoteAddr)
		} else {
			log.Printf("error while registering (create) client player : '%v'", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		name, err := ParseSingleStringUrlQueryParam(request, "name")
		emotar, err := ParseSingleStringUrlQueryParam(request, "emotar")
		player.Name = *name
		player.Emotar = *emotar
		_, err = playersRepository.UpdatePlayer(*player)
		if err != nil {
			log.Printf("error while registering (update) client player : '%v'", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Using existing (and updated) player %+v for ip=%s \n", player, request.RemoteAddr)
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func createPlayer(id int, name string, emotar string) (*repositories.Player, error) {
	return &repositories.Player{
		Name:   name,
		Id:     id,
		Emotar: emotar,
	}, nil
}

func GetPlayerById(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		log.Printf("error while retrieving player(id=%d): '%v'\n", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func UpdatePlayer(response http.ResponseWriter, request *http.Request) {
	var player repositories.Player
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := playersRepository.UpdatePlayer(player)
	if err != nil {
		log.Printf("error while updating Player: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeletePlayer(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = playersRepository.DeletePlayer(id)
	if err != nil {
		log.Printf("error while deleting player(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
