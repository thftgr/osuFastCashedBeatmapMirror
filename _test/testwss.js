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


function lb() {
    let gvar = {
        host: "api.nerina.pw", latency: null
    }
    const apiserver = [
        "api.nerinyan.moe",
        "api.nerinyan.moe",
        "rus.nerinyan.moe",
        "rus.nerinyan.moe",
        "ko2.nerinyan.moe",
        "ko2.nerinyan.moe",
        "api.nerinyan.moe",
        "api.nerinyan.moe",
        "rus.nerinyan.moe",
        "rus.nerinyan.moe",
        "ko2.nerinyan.moe",
        "ko2.nerinyan.moe",

    ]

    gvar.latency = null

    const st = new Date().getTime()

    apiserver.map((host) => {
        const promises = [];
        promises.push(fetch(`https://${host}/health`).then(res => {
            if (!res.ok) throw null
        }));

        Promise.all(promises).then(() => {

            const et = new Date().getTime()
            if (gvar.latency === null || gvar.latency > (et - st)){
                gvar.host = host
                gvar.latency = et - st
            }

            console.log(`host: ${host} ${et - st} ms`) // 이
        }).catch(e => console.log(`host: ${host} error : ${e}`));

    })
}

lb();