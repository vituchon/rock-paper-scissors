var app = {
  game: null,
  clientPlayer: null,
  soundAllowed: false,
  scoreByPlayerId: {},
  getEnemyPlayer() {
    return this.game.players.find(player => player.id !== this.clientPlayer.id);
  },
  initScores() {
    if (this.game.players.length !== 2) {
      throw new Error("invalid score initialization, it requires to have two players");
    }
    this.game.players.forEach((player) => {
      this.scoreByPlayerId[player.id] = 0
    })
  }
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
    const timeouts = ids.map(id => {
      return new Promise((resolve) => {
        let button = document.getElementById(id);
        if (button) {
          button.disabled = true;
          button.classList.remove('fade-in');
          button.classList.add('fade-out');
          setTimeout(() => {
            button.style.display = "none";
            resolve();
          }, 1000);
        } else {
          resolve();
        }
      });
    });

    return Promise.all(timeouts);
  }

  function fadein(display, ...ids) {
    const timeouts = ids.map(id => {
      return new Promise((resolve) => {
        let button = document.getElementById(id);
        if (button) {
          button.style.display = display;
          button.classList.remove('fade-out');
          button.classList.add('fade-in');
          setTimeout(() => {
            button.disabled = false;
            resolve();
          }, 1000);
        } else {
          resolve();
        }
      });
    });

    return Promise.all(timeouts);
  }

  Buttons.fadeout = fadeout
  Buttons.fadein = fadein
})(Buttons || (Buttons = {}));


function displayMessage(message) {
  Toastify({
    text: message,
    duration: 3000,
    newWindow: true,
    gravity: "top",
    position: 'center',
    style: {
      "font-size": "1.5em",
      "border-radius": "0.5em",
    }
  }).showToast();
}


function displayErrorMessage(message) {
  Toastify({
    text: message,
    duration: 3000,
    newWindow: true,
    gravity: "bottom",
    position: 'left',
    style: {
      "font-size": "1em",
      "border-radius": "0.5em",
      "background": "linear-gradient(135deg, #f0000f, #5477f5)",
      "color": "yellow",
    }
  }).showToast();
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