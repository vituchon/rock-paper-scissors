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
