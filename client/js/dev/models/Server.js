define([
    "jquery",
    "util/socket",
    "Backbone",
    "Underscore"
], function($, socket, Backbone, _) {

    var Server = Backbone.Model.extend({
        idAttribute: "url",

        defaults: {
            "state": "disconnected",
            "dConnected": $.Deferred()
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
            this.set("state", "connected");
            this.get("dConnected").resolve();
        },

        onMessage: function(evt) {
            console.log("Got message!", evt.data);
        },
     
        onClose: function() {
            if (this.get("state") != "disconnecting") {
                console.log("Connection died", this.toString());
            }
            this.set("state", "disconnected");
            console.log("Disconnected", this.toString());
        },

        toString: function() {
            return this.get("url") + " [" + this.get("state") + "]";
        }

    });

    return Server;
});
