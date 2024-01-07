function (cookie) {
    var url = new URL(document.currentScript.src);
    url.protocol = url.protocol == "https:" ? "wss:" : "ws:";

    var last = cookie;
    var timeout;

    function resetBackoff() {
        timeout = 100;
    };

    function backoff() {
        if (timeout > 2000) {
            return;
        }
        timeout = timeout * 2;
    };

    function connect() {
        const socket = new WebSocket(url);

        socket.onmessage = (event) => {
            if (event.data !== "") {
                if (last === "") {
                    last = event.data;
                }

                if (last !== event.data) {
                    console.log("[Hot Reload] Reloading");
                    socket.close();
                    location.reload();
                }
            }

            socket.send("");
        };

        socket.onopen = () => {
            resetBackoff();
            console.log("[Hot Reload] Connected");
        };

        socket.onclose = () => {
            const id = setTimeout(function () {
                clearTimeout(id);
                backoff();
                connect();
            }, timeout);
        };
    }

    resetBackoff();
    connect();
}
