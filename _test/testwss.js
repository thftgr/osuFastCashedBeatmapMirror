let socket = new WebSocket("wss://wss.thftgr.synology.me/ws");

socket.onopen = function(e) {
    console.log(`${new Date().toISOString()} [connected] websocket server:`);
    socket.send("ping");
    socket.send("ping");
    socket.send("ping");
    socket.send("ping");
    socket.send("ping");
};

socket.onmessage = function(event) {
    console.log(`${new Date().toISOString()} [message] from server: ${JSON.parse(event.data)[0].id}`);
};

socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`[close] connection closed (code=${event.code} reason=${event.reason})`);
    } else {
        console.log('[close] connection dead.');
    }
};

socket.onerror = function(error) {
    console.log(`[error] ${error.message}`);
};