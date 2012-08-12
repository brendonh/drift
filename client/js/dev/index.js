require.config({
    "paths": {
        "jquery": "require-jquery",
        "text": "libs/requirejs/text",
        "dispatch": "libs/dispatch",
        "router": "libs/router",
        "Underscore": "libs/underscore-min",
        "Backbone": "libs/backbone-min",
        "msgpack": "libs/msgpack"
    },
    "shim": {
        "Underscore": {
            "deps": ["jquery"],
            "exports": "_"
        },
        "Backbone": {
            "deps": ["Underscore", "jquery"],
            "exports": "Backbone"
        }
    },
    "deps": ["index"],
    "baseUrl": "/js/dev"
})

require([
    "jquery",
    "Backbone",
    "collections/Servers"
], function($, Backbone, Servers) {

    var initialServer = "ws://dev.brendonh.org:9998/";

    var servers = new Servers();

    window.onbeforeunload = function() {
        servers.closeAll();
    }

    servers.ensure(initialServer).done(startLogin);
});

function startLogin(server) {
    console.log("Connected!");
    
    var user = getQueryVariable('user');
    var password = getQueryVariable('password');
    
    if (!user || !password) {
        alert("Give user and password in query string");
        return;
    }
    
    server.callAPI("accounts", "login", 
                   {"name": user,
                    "password": password})
        .done(function(r) { afterLogin(server, r); });
}

function afterLogin(server, response) {
    var shipID = getQueryVariable('ship');
    if (!shipID) {
        server.callAPI("ships" ,"list")
            .done(function(shipResponse) {
                var ships = shipResponse["data"]["ships"];
                for (var i in ships) {
                    var ship = ships[i];
                    console.log(ship.name, ship.id);
                }
            });
        return;
    }

    server.callAPI("ships", "control", {"id": shipID})
        .done(function() {
            console.log("Ready!");
        });
}
    

function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split('&');
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) == variable) {
            return decodeURIComponent(pair[1]);
        }
    }
    return undefined;
}