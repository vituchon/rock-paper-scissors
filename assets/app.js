var app = {
  game: null,
  clientPlayer: null,
  soundAllowed: false,
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
  } else {
    console.warn("event",event,"is already registered for element", element, ". Skipping registration.")
  }
}


var Buttons;
(function (Buttons) {

  function fadeout(...ids) {
    ids.forEach(id => {
      let button = document.getElementById(id);
      if (button) {
        button.disabled = true;
        button.classList.remove('fade-in');
        button.classList.add('fade-out');
        setTimeout(() => {
          button.style.display = "none";
        }, 1000);
      }
    });
  }

  function fadein(display, ...ids) {
    ids.forEach(id => {
      let button = document.getElementById(id);
      if (button) {
        button.style.display = display;
        button.classList.remove('fade-out');
        button.classList.add('fade-in');
        setTimeout(() => {
          button.disabled = false;
        }, 1000);
      }
    });
  }

  Buttons.fadeout = fadeout
  Buttons.fadein = fadein
})(Buttons || (Buttons = {}));

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