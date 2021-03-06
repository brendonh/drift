require.config({
    "paths": {
        "jquery": "require-jquery",
        "text": "libs/requirejs/text",
        "dispatch": "libs/dispatch",
        "router": "libs/router",
        "Underscore": "libs/underscore-min",
        "Backbone": "libs/backbone-min",
        "msgpack": "libs/msgpack",
        "Three": "libs/Three"
    },
    "shim": {
        "Underscore": {
            "deps": ["jquery"],
            "exports": "_"
        },
        "Backbone": {
            "deps": ["Underscore", "jquery"],
            "exports": "Backbone"
        },
        "Three": {
            "exports": "THREE"
        }
    },
    "deps": ["index"],
    "baseUrl": "/js/dev"
})

require([
    "jquery",
    "Backbone",
    "collections/Servers",
    "models/Sector",
    "models/SectorRenderer",
    "views/SectorRendererView",
    "models/Ship"
], function($, Backbone, Servers, Sector, SectorRenderer, SectorRendererView, Ship) {

    var initialServer = "ws://dev.brendonh.org:9998/";

    var servers = new Servers();

    window.onbeforeunload = function() {
        servers.closeAll();
    }

    servers.ensure(initialServer).done(startLogin);

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
        if (shipID) {
            chooseShip(server, shipID);
            return;
        }
         
        server.callAPI("ships" ,"list")
            .done(function(shipResponse) {
                var ships = shipResponse["data"]["ships"];
                if (!ships.length) {
                    console.log("No ships :(");
                    server.callAPI("ships", "create", 
                                   {name: "Fluffy"}
                                  ).done(function(newShipResponse) {
                                      var id = newShipResponse["data"]["id"];
                                      console.log("Created ship", id)
                                      chooseShip(server, id);
                                  });
                    return;
                }
                for (var i in ships) {
                    var ship = ships[i];
                    console.log(ship.name, ship.id);
                }
                chooseShip(server, ships[0].id);
            });
    }

    function chooseShip(server, shipID) {
        server.callAPI("server", "control", {
            "id": shipID
        }).done(function() {
            afterShip(server, shipID);
        });
    }

    function afterShip(server, shipID) {
        var sector = new Sector();

        var renderer = new SectorRenderer({
            "sector": sector
        });

        var view = new SectorRendererView(
            {"el": $("#sector"), 
             "model": renderer});
        view.render();

        var ship = new Ship({"id": shipID})
        sector.get("ships").add(ship);

        var ghost = new Ship({"id": "ghost"})
        sector.get("ships").add(ghost);
        ghost.thrust(true);
        ghost.rotateRight(true);

        window.onkeydown = function(e) {
            if (e.keyCode == 38) ship.thrust(true);
            else if (e.keyCode == 39) ship.rotateRight(true);
            else if (e.keyCode == 37) ship.rotateLeft(true);
        };
        
        window.onkeyup = function(e) {
            if (e.keyCode == 38) ship.thrust(false);
            else if (e.keyCode == 39) ship.rotateRight(false);
            else if (e.keyCode == 37) ship.rotateLeft(false);
        };

        var requestAnimFrame = (function(){
            return  window.requestAnimationFrame   || 
                window.webkitRequestAnimationFrame || 
                window.mozRequestAnimationFrame    || 
                window.oRequestAnimationFrame      || 
                window.msRequestAnimationFrame     || 
                function( callback ){
                    window.setTimeout(callback, 1000 / 60);
                };
        })();

        var frames = 0;
        var lastInfoTime = new Date().getTime();

        function frame() {
            sector.tick();
            view.renderFrame();
            info();
            requestAnimFrame(frame);
        }

        function info() {
            if (frames++ > 60) {
                var time = new Date().getTime();
                var delta = (time - lastInfoTime) / 1000;
                
                $("#info .fps").html( Math.round((frames / delta)) );
                
                frames = 0;
                lastInfoTime = time;

                var faces = 0;
                sector.get("ships").each(function(ship) {
                    faces += ship.get("mesh").geometry.faces.length;
                });
                    
                $("#info .faces").html(faces);
            }
        }

        frame(frame);
    }

});


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