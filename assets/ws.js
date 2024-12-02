var WebSockets;
(function (WebSockets) {
    const normalCloseEventCode = 1000;
    const closeDescriptionByCode = {
        [normalCloseEventCode]: "CLOSE_NORMAL",
        1001: "CLOSE_GOING_AWAY",
        1002: "CLOSE_PROTOCOL_ERROR",
        1003: "CLOSE_UNSUPPORTED",
        1004: "Reserved",
        1005: "CLOSED_NO_STATUS",
        1006: "CLOSE_ABNORMAL",
        1007: "Unsupported payload",
        1008: "Policy violation",
        1009: "CLOSE_TOO_LARGE",
        1010: "Mandatory extension",
        1011: "Server error",
        1012: "Service restart",
        1013: "Try again later",
        1014: "Bad gateway",
        1015: "TLS handshake fail"
    };

    function resolveProtocol() {
        const isSecure = window.location.protocol.indexOf("https") !== -1;
        return isSecure ? "wss" : "ws";
    }

    function resolveHost() {
        return window.location.host;
    }

    class Service {

        retrieve() {
            if (typeof this.webSocket === "undefined") {
                return this.adquire();
            } else {
                if (this.webSocket.readyState === WebSocket.OPEN) {
                    console.debug("Web socket already acquired and open");
                    return Promise.resolve(this.webSocket);
                } else {
                    return Promise.reject("already bound to another tab");
                }
            }
        }

        adquire() {
            const deferred = new Promise((resolve, reject) => {
                try {
                    const protocol = resolveProtocol();
                    const host = resolveHost();
                    this.webSocket = new WebSocket(`${protocol}://${host}/adquire-ws`);

                    this.webSocket.onopen = (event) => {
                        window.addEventListener("beforeunload", () => {
                            cleanup();
                        });
                        resolve(this.webSocket);
                    };

                    this.webSocket.onerror = (event) => {
                        console.warn("Web socket error, event is: ", event);
                        cleanup();
                        reject("error on websocket 😢");
                    };

                    this.webSocket.onclose = (event) => {
                        const reason = event.reason || closeDescriptionByCode[event.code];
                        console.debug("Closing web socket. Was clean:", event.wasClean, ", code:", event.code, ", reason:", reason);
                        const hasNotCloseNormal = event.code !== normalCloseEventCode;
                        if (hasNotCloseNormal || !event.wasClean) {
                            reject(reason);
                        }
                        this.webSocket?.removeEventListener?.("message", handleServerMessage);
                        this.webSocket = undefined;
                    };

                    const handleServerMessage = (event) => {
                        const notification = JSON.parse(event.data);
                        if (notification.kind === "broadcast") {
                            Toastify({
                              text: notification.message,
                              duration: 4500,
                              newWindow: true,
                              gravity: "top",
                              position: 'center',
                            }).showToast();
                        }
                    };

                    const cleanup = () => {
                        this.webSocket?.removeEventListener?.("message", handleServerMessage);
                        this.release();
                    };

                    this.webSocket.addEventListener("message", handleServerMessage);
                } catch (error) {
                    reject(error);
                }
            });

            return deferred;
        }

        release(code, reason) {
            if (typeof this.webSocket === "undefined") {
                return Promise.reject("No web socket to release");
            }
            this.webSocket.close(code, reason); // honoring convention (but it is not necessary)
            return fetch("/release-ws")
                .then(() => {})
                .catch((error) => console.error("Release WebSocket error: ", error));
        }

        pingme(message) {
            const config = {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                params: { message }
            };

            const url = `/send-message-ws?message=${encodeURIComponent(message)}`;
            return fetch(url, config)
                .then(response => response.json())
                .catch((error) => console.error("Pingme error: ", error));
        }

        pingall(message) {
            const config = {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                params: { message }
            };

            const url = `/send-message-all-ws?message=${encodeURIComponent(message)}`;
            return fetch(url, config)
                .then(response => response.json())
                .catch((error) => console.error("Pingall error: ", error));
        }
    }

    WebSockets.DefaultService = new Service()
})(WebSockets || (WebSockets = {}));
