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

        sendAPICall(login, socket)
        sendAPICall(ping, socket)

    }

    socket.ondata = function() {
        console.log("DATA")
    }

    socket.onmessage = function(evt) {
        var data = evt.data;

        // XXX BGH TODO: Use a better msgpack implementation
        var bytes = []
        var bufView = new Uint8Array(data)
        for (var i = 0; i < data.byteLength; i++) {
            bytes.push(bufView[i]);
        }
        console.log("Message:", JSON.stringify(msgpack.unpack(bytes)));
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
    console.log("Call:", JSON.stringify(data))
    var msg = msgpack.pack(data);

    // XXX BGH TODO: Use a better msgpack implementation
    var buf = new ArrayBuffer(msg.length + 1);
    var bufView = new Uint8Array(buf);
    bufView[0] = 'a'.charCodeAt(0);
    for (var i=0; i<msg.length; i++) {
        bufView[i+1] = msg[i];
    }
 
    socket.send(buf);
}