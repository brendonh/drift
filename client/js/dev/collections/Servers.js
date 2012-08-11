define([
    "Backbone",
    "models/Server"
], function(Backbone, Server) {

    var Servers = Backbone.Collection.extend({
        model: Server,

        ensure: function(url) {
            var server = this.get(url);
            if (!server) {
                server = new Server({url: url})
                this.add(server);
            }
            server.connect();
            return server.get("dConnected");
        },

        closeAll: function() {
            this.each(function(socket) {
                socket.disconnect();
            });
        }
        
    });

    return Servers;

});