package controllers

import (
	"net/http"
	"log"

	"github.com/gorilla/sessions"
	"github.com/vituchon/escobita/util"
)

var clientSessions *sessions.CookieStore
var integerSequence util.IntegerSequence

func InitSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
	clientSessions.Options = &sessions.Options{
		HttpOnly: true,         // La cookie solo será accesible por HTTP
		SameSite: http.SameSiteLaxMode, // Establece el comportamiento de SameSite
		MaxAge:   3600,         // Establece el tiempo de vida de la cookie (en segundos)
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
	log.Printf("clientSession: '%+v' /n",clientSession)
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
