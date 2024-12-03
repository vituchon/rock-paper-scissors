package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/vituchon/escobita/util"

	"errors"
	"fmt"

	"github.com/vituchon/rock-paper-scissors/repositories"
)

var clientSessions *sessions.CookieStore
var integerSequence util.IntegerSequence

func InitSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
	clientSessions.Options = &sessions.Options{
		HttpOnly: true,                 // La cookie solo será accesible por HTTP
		SameSite: http.SameSiteLaxMode, // Establece el comportamiento de SameSite
		MaxAge:   3600,                 // Establece el tiempo de vida de la cookie (en segundos)
	}
	integerSequence = util.NewFsIntegerSequence("ppt.seq", 0, 1)
}

func GetOrCreateClientSession(request *http.Request) (*sessions.Session, error) {
	clientSession, err := clientSessions.Get(request, "ppt_client")
	if err != nil {
		return nil, err
	}
	if clientSession.IsNew {
		nextId, err := integerSequence.GetNext()
		if err != nil {
			return nil, err
		}
		clientSession.Values["clientId"] = nextId
	}
	log.Printf("clientSession: '%+v'", clientSession)
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}

var ClientPlayerDoesntExistsErr = errors.New("Client player doesn't exists")

func GetClientPlayer(request *http.Request) (*repositories.Player, error) {
	clientSession := request.Context().Value("clientSession").(*sessions.Session)
	wrappedInt, exists := clientSession.Values["clientId"]
	if !exists {
		return nil, ClientPlayerDoesntExistsErr
	}

	id, ok := wrappedInt.(int)
	if !ok {
		return nil, fmt.Errorf("Client id exists and is '%v' but it can not be asserted as an int", wrappedInt)
	}
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		err := fmt.Errorf("error while retrieving player(id='%d'): '%v'", id, err)
		return nil, err
	}
	return player, nil
}
