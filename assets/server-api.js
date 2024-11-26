"use strict";
function handleFetchResponseError(response) {
    return response.text().then((text) => {
        text = text || "Error";
        throw new Error(`${text}, status code: ${response.status}`);
    });
}
var ServerApi;
(function (ServerApi) {
    class Service {
        constructor() { }

        registerPlayer(player) {
            const requestInit = {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
            };
            return fetch("/api/v1/players/register?name="+player.name+"&emotar="+player.emotar, requestInit).then((r) => {
                if (!r.ok) {
                    return handleFetchResponseError(r);
                }
                return r.json().then((loggedPlayer) => {
                    return loggedPlayer;
                });
            }).catch((err) => {
                console.error('Failed to login:', err);
                throw err;
            });
        }

        joinGame(player, gameId) {
            const requestInit = {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(player)
            };
            return fetch(`/api/v1/join`, requestInit).then((r) => {
                if (!r.ok) {
                    return handleFetchResponseError(r);
                }
                return r.json().then((game) => {
                    return game;
                });
            }).catch((err) => {
                console.error('Failed to join game:', err);
                throw err;
            });
        }
    }
    ServerApi.Service = Service;
    ServerApi.DefaultService = new Service();
})(ServerApi || (ServerApi = {}));
