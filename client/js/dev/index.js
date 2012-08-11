require.config({
    "paths": {
        "jquery": "require-jquery",
        "text": "libs/requirejs/text",
        "dispatch": "libs/dispatch",
        "router": "libs/router",
        "Underscore": "libs/underscore-min",
        "Backbone": "libs/backbone-min"
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
], function($, Backbone, Servers){
    var servers = new Servers();

    window.onbeforeunload = function() {
        servers.closeAll();
    }

    var d = servers.ensure("ws://dev.brendonh.org:9998/")
    d.done(function() {
        console.log("Connected!");
    });
});
