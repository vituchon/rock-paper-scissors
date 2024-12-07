"use strict";
function handleFetchResponseError(response) {
  return response.text().then((text) => {
      text = text || "Error";
      throw new Error(`${text}, status code: ${response.status}`);
  });
}

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

  getGames() {
      const requestInit = {
          method: 'GET',
          headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
          },
      };
      return fetch(`/api/v1/games`, requestInit).then((r) => {
          if (!r.ok) {
              return handleFetchResponseError(r);
          }
          return r.json().then((games) => {
              return games;
          });
      }).catch((err) => {
          console.error('Failed to fetch games:', err);
          throw err;
      });
  }

  createGame(game) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(game)
    };
    return fetch(`/api/v1/games`, requestInit).then((r) => {
        if (!r.ok) {
            return handleFetchResponseError(r);
        }
        return r.json().then((games) => {
            return games;
        });
    }).catch((err) => {
        console.error('Failed to fetch games:', err);
        throw err;
    });
  }

  joinGame(gameId) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    };
    return fetch(`/api/v1/games/${gameId}/join`, requestInit).then((r) => {
        if (!r.ok) {
            return handleFetchResponseError(r);
        }
        return r.json().then((game) => {
            return this.bindWebSocket(game.id).then(() => {
              return game
            })
        });
    }).catch((err) => {
        console.error('Failed to join game:', err);
        throw err;
    });
  }
  bindWebSocket(id) {
    return fetch(`/api/v1/games/${id}/bind-ws`)
  }

  quitGame(game) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    };
    return fetch(`/api/v1/games/${game.id}/quit`, requestInit).then((r) => {
      if (!r.ok) {
            return handleFetchResponseError(r);
        }
        return r.json().then((game) => {
            return this.unbindWebSocket(game.id).then(() => {
              return game
            })
        });
    }).catch((err) => {
        console.error('Failed to join game:', err);
        throw err;
    });
  }
  unbindWebSocket(id) {
    return fetch(`/api/v1/games/${id}/unbind-ws`)
  }

  startGame(game) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    };
    return fetch(`/api/v1/games/${game.id}/start`, requestInit).then((r) => {
        if (!r.ok) {
          return handleFetchResponseError(r);
        }
        return r.json().then((game) => {
          return game
        });
    }).catch((err) => {
        console.error('Failed to join game:', err);
        throw err;
    });
  }

  restartGame(game) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    };
    return fetch(`/api/v1/games/${game.id}/restart`, requestInit).then((r) => {
        if (!r.ok) {
          return handleFetchResponseError(r);
        }
        return r.json().then((game) => {
          return game
        });
    }).catch((err) => {
        console.error('Failed to restart game:', err);
        throw err;
    });
  }

  performAction(game, action) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(action)
    };
    return fetch(`/api/v1/games/${game.id}/perform-action`, requestInit).then((r) => {
        if (!r.ok) {
          return handleFetchResponseError(r);
        }
        return r.json().then((game) => {
          return game
        });
    }).catch((err) => {
        console.error('Failed to perform action:', err);
        throw err;
    });
  }

  resolveCurrentMatch(game) {
    const requestInit = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        }
    };
    return fetch(`/api/v1/games/${game.id}/resolve-current-match`, requestInit).then((r) => {
      if (!r.ok) {
        return handleFetchResponseError(r);
      }
      return r.json().then((winnerPlayerId) => {
        return winnerPlayerId
      });
    }).catch((err) => {
      console.error('Failed to resolve current match:', err);
      throw err;
    });
  }

  deleteAllGames() {
    const requestInit = {
        method: 'DELETE',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    };
    return fetch(`/api/v1/games`, requestInit).then((r) => {
      if (!r.ok) {
        return handleFetchResponseError(r);
      }
      return r.ok
    }).catch((err) => {
      console.error('Failed to resolve current match:', err);
      throw err;
    });
  }
}

const DefaultService = new Service();

export { DefaultService, Service };