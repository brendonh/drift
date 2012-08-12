define([
    "jquery",
    "Backbone",
    "Underscore",
    "collections/SectorShips"
], function($, Backbone, _, SectorShips) {
    var Sector = Backbone.Model.extend({
        idAttribute: "coordString",

        initialize: function() {
            this.set("ships", new SectorShips());
        },

        tick: function() {
            this.get("ships").each(function(ship) {
                ship.tick();
            });
            this.trigger("tick");
        }

    });
    
    return Sector;
});