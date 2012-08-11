function startSocket(url) {
    var socket;

    if ("WebSocket" in window) {
        socket = new WebSocket(url);
    } else if ("MozWebSocket" in window) {
        socket = new MozWebSocket(url);
    }

    if (!socket) {
        alert("No websocket support!");
        return;
    }

    socket.binaryType = "arraybuffer";

    socket.onopen = function() {
        console.log("Connected!");

        var ping = {
            'service': 'accounts',
            'method': 'ping',
            'data': {}
        }

        var login = {
            'service': 'accounts',
            'method': 'login',
            'data': {
                'name': 'other',
                'password': 'test'
            }
        }

        var ships = {
            'service': 'ships',
            'method': 'list',
            'data': {}
        }

        sendAPICall(login, socket)

        // Temp until we have callbacks
        setTimeout(
            function() {
                // sendAPICall({"service": "ships",
                //              "method": "create",
                //              "data": { "name": "Sparky" }},
                //            socket)
                //sendAPICall(ships, socket)
                sendAPICall(
                    {"service": "ships",
                     "method": "control",
                     "data": { "id": "d4e15abf-bf62-4bae-bf88-ce63026ded41"}},
                    socket)
            }, 500)

    }

    socket.ondata = function() {
        console.log("DATA")
    }

    socket.onmessage = function(evt) {
        var buf = evt.data;

        var view = new DataView(buf);
        var msgType = String.fromCharCode(view.getUint8(0));

        if (msgType == 'a') {
            var response = msgpack.decodeFromView(view, 1);
            console.log("API reply:", JSON.stringify(response))
        }
    }

    socket.onclose = function() {
        console.log("Disconnected!");
    }

    window.onbeforeunload = function() {
        socket.onclose = function () {};
        socket.close();
    };

}

var apiID = 0;

function sendAPICall(data, socket) {
    data['id'] = apiID++;
    console.log("Call:", JSON.stringify(data));

    var buf = new ArrayBuffer(64 * 1024);
    var view = new DataView(buf);
    var len = msgpack.encodeToView(data, view, 1);
    view.setUint8(0, 'a'.charCodeAt(0));

    socket.send(buf.slice(0, len+1));
}