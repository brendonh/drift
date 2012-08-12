define([
    "Backbone",
    "models/Ship"
], function(Backbone, Ship) {

    var SectorShips = Backbone.Collection.extend({
        model: Ship
        
    });

    return SectorShips;

});