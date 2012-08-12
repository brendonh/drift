define([
    "Backbone",
    "models/CubeSet"
], function(Backbone, CubeSet) {

    var Ship = Backbone.Model.extend({
        
        defaults: {
            "moving": false,
            "rotateLeft": false,
            "rotateRight": false,

            "flip": 0,
            "rotation": 0
        },

        initialize: function() {

            var colors = {
                'X': 0x993333,
                'x': 0x999999,
                'o': 0x444444
            };

            var cubes = new CubeSet(10);
            cubes.addAscii([
                ' o o ',
                ' o o ',
                ' oxo ',
                ' xxx ',
                'xxXxx',
                'xxxxx',
                ' xox '
            ], colors);

            this.set("cubes", cubes);

            this.set("mesh", cubes.toMesh());
        },

        thrust: function(onOff) {
            this.set("moving", onOff);
        },

        rotateRight: function(onOff) {
            this.set("rotateRight", onOff);
        },

        rotateLeft: function(onOff) {
            this.set("rotateLeft", onOff);
        },
        
        tick: function() {
            var state = this.attributes;

            if (state['rotateLeft']) {
                state['flip'] -= (state['flip'] + 0.8) * 0.1;
                state['rotation'] -= 0.05;
            } else if (state['rotateRight']) {
                state['flip'] += (0.8 - state['flip']) * 0.1;
                state['rotation'] += 0.05;
            } else {
                state['flip'] *= 0.9;
                if (Math.abs(state['flip']) < 0.1) state['flip'] = 0;
            }

            var speed = 10;

            if (state['moving']) {
                state['mesh'].position.x -= speed * Math.sin(state['rotation']);
                state['mesh'].position.y += speed * Math.cos(state['rotation']);
            }

            state['mesh'].quaternion.setFromAxisAngle(
                new THREE.Vector3(0, 0, 1), state['rotation']);
            var flip = new THREE.Quaternion();
            flip.setFromAxisAngle(new THREE.Vector3(0, 1, 0), state['flip']);
            state['mesh'].quaternion.multiplySelf(flip);
        }
    });
                                     
    return Ship;

});