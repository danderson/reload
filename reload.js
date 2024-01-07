function (cookie) {
    var url = new URL(document.currentScript.src);
    url.protocol = url.protocol == "https:" ? "wss:" : "ws:";

    var retry = 50; // ms

    function connect() {
        retry = Math.min(retry*2, 2000);

        const socket = new WebSocket(url);

        socket.onmessage = (event) => {
            if (event.data !== "") {
                if (cookie !== event.data) {
                    console.log("[Hot Reload] Reloading");
                    socket.close();
                    location.reload();
                }
            }

            socket.send("");
        };

        socket.onopen = () => {
            retry = 100;
            console.log("[Hot Reload] Connected");
        };

        socket.onclose = () => {
            const id = setTimeout(function () {
                clearTimeout(id);
                connect();
            }, retry);
        };
    }

    connect();
}
