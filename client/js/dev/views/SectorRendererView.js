define([
    "jquery",
    "Backbone",
    "Three"
], function($, Backbone, THREE) {

    var SectorRendererView = Backbone.View.extend({

        initialize: function() {
        },
        
        render: function() {
            var renderer = this.model.get("renderer");

            $(renderer.domElement)
                .css("height", this.model.get("height") + "px")
                .css("width", this.model.get("width") + "px")
                .css("background-color", "black");

            this.$el.empty().append(renderer.domElement);
        },

        renderFrame: function() {
            this.model.get('renderer').render( 
                this.model.get('scene'), 
                this.model.get('camera'));
        }

    });

    return SectorRendererView;

});
    