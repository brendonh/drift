define([
    "jquery",
    "util/socket",
    "Backbone",
    "Underscore",
    "libs/msgpack"
], function($, socket, Backbone, _, msgpack) {

    var Server = Backbone.Model.extend({
        idAttribute: "url",

        defaults: {
            "state": "disconnected",
            "dConnected": $.Deferred(),

            "apiID": 0,
        },

        connect: function() {
            if (this.get("state") != "disconnected") {
                console.log("Already connecting", this.toString());
                return;
            }

            var url = this.get("url");
            console.log("Connecting to " + url);

            socket = new socket(url);
            socket.binaryType = "arraybuffer";
            socket.onopen = _.bind(this.onConnect, this);
            socket.onmessage = _.bind(this.onMessage, this);
            socket.onclose = _.bind(this.onClose, this);

            this.set("state", "connecting");
            this.set("socket", socket);
        },
        
        disconnect: function() {
            if (this.get("state") == "new") {
                console.log("Not connected", this.toString);
                return;
            }

            this.set("state", "disconnecting");                
            console.log("Disconnecting", this.toString());

            this.get("socket").close();
        },

        onConnect: function() {
            console.log("Connected", this.toString());
            this.set("callbacks", {});
            this.set("state", "connected");
            this.get("dConnected").resolve(this);
        },

        onMessage: function(evt) {
            var buf = evt.data;

            var view = new DataView(buf);
            var msgType = String.fromCharCode(view.getUint8(0));
            
            if (msgType == 'a') {
                var response = msgpack.decodeFromView(view, 1);
                console.log("API reply:", JSON.stringify(response))
                var id = response['id'];
                var callbacks = this.get("callbacks");
                var d = callbacks[id];
                if (d) {
                    delete callbacks[id];
                    d.resolve(response);
                } else {
                    console.log("Reply to unknown API call:", response);
                }
            } else {
                console.log("Unknown packet:", dataViewToString(view));
            }
        },
     
        onClose: function() {
            if (this.get("state") != "disconnecting") {
                console.log("Connection died", this.toString());
            }
            this.set("state", "disconnected");
            console.log("Disconnected", this.toString());
        },

        callAPI: function(service, method, data) {
            if (this.get("state") != "connected") {
                console.log("Unconnected server ignoring API call", 
                            service, method, data, this.toString());
                return;
            }

            var id = this.get("apiID");
            this.set("apiID", id + 1);

            var request = {
                "id": id,
                "service": service,
                "method": method,
                "data": data || {}
            };

            var buf = new ArrayBuffer(64 * 1024);
            var view = new DataView(buf);
            var len = msgpack.encodeToView(request, view, 1);
            view.setUint8(0, 'a'.charCodeAt(0));

            this.get("socket").send(buf.slice(0, len+1));

            var d = $.Deferred();
            this.get("callbacks")[id] = d;
            return d;
        },


        toString: function() {
            return this.get("url") + " [" + this.get("state") + "]";
        }

    });

    function dataViewToString(view) {
        var out = "";
        for (var i = 0; i < view.byteLength; i++) {
            out += String.fromCharCode(view.getUint8(i));
        }
        return out;
    }

    return Server;
});
