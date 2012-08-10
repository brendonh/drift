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

        var login = {
            'service': 'accounts',
            'method': 'login',
            'data': {
                'name': 'brendonh',
                'password': 'test'
            }
        }

        var ping = {
            'service': 'accounts',
            'method': 'ping',
            'data': {}
        }

        sendAPICall(ping, socket)
        sendAPICall(login, socket)
        sendAPICall(ping, socket)

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