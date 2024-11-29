var app = {
  game: null,
  clientPlayer: null,
  soundAllowed: true,
  scoreByPlayerId: {},
  getEnemyPlayer() {
    return this.game.players.find(player => player.id !== this.clientPlayer.id);
  },
};

function split(array, predicate) {
  const results = array.reduce(([passed, rejected], elem) => {
    if (predicate(elem)) {
      passed.push(elem);
    } else {
      rejected.push(elem)
    }
    return [passed, rejected]
  }, [[], []]);
  return results
}

function addEventIfNotRegistered(element, event, callback) {
  const attributeName = `data-event-${event}`;
  const isEventRegistered = element.getAttribute(attributeName);

  if (!isEventRegistered || isEventRegistered === "false") {
      element.addEventListener(event, callback);
      element.setAttribute(attributeName, "true");
  }
}

function watch (oObj, sProp) {
  var sPrivateProp = "$_"+sProp+"_$"; // to minimize the name clash risk
  oObj[sPrivateProp] = oObj[sProp];

  // overwrite with accessor
  Object.defineProperty(oObj, sProp, {
      get: function () {
          return oObj[sPrivateProp];
      },

      set: function (value) {
          //console.log("setting " + sProp + " to " + value);
          debugger; // sets breakpoint
          oObj[sPrivateProp] = value;
      }
  });
}