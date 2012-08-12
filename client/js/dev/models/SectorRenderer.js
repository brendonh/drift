define([
    "Backbone",
    "Three"
], function(Backbone, THREE) {

    var SectorRenderer = Backbone.Model.extend({
        
        defaults: {
            "width": 800,
            "height": 600,
        },

        initialize: function() {
            var scene = new THREE.Scene();

            var camera = new THREE.PerspectiveCamera( 60, 800 / 600, 1, 2000 );
            camera.position.z = 1000;
            scene.add(camera);

            var renderer = new THREE.CanvasRenderer();
            renderer.setSize(this.get("width"), this.get("height"));

            var directionalLight = new THREE.DirectionalLight( 0xffffff, 1.0 );
	        directionalLight.position.set( 0.5, 0.2, -1 );
            scene.add(directionalLight);

            this.set({
                "scene": scene,
                "camera": camera,
                "renderer": renderer
            });

            this.addDust();

            this.get("sector").get("ships").on("add", this.addShip, this);

            this.get("sector").on("tick", this.tick, this);
            
        },

        addShip: function(ship) {
            this.get("scene").add( ship.get("mesh") );
        },
        
        addDust: function() {
            
            var particleRender = function(context) { 
	            context.beginPath();
	            context.arc( 0, 0, 1, 0,  Math.PI * 2, true );
	            context.fill();
            };
            
            var count = 200,
            particles = [],
            pMaterial =
                new THREE.ParticleCanvasMaterial( { 
                    color: 0xffffff, 
                    program: particleRender } );

            var scene = this.get("scene");
            
            for(var p = 0; p < count; p++) {
                var particle = new THREE.Particle(pMaterial);
                
                particle.position.x = -1000 + Math.random() * 2000;
                particle.position.y = -1000 + Math.random() * 2000;
                particle.position.z = Math.random() * 500;
                
                particles.push(particle);
                scene.add(particle);
            }
            
            this.set("dust", particles);
        },

        tick: function() {
            var state = this.attributes;

            var watchedShip = this.get("sector").get("ships").at(0).get("mesh");

            state['camera'].position = watchedShip.position.clone();
            state['camera'].position.z = -700;
            state['camera'].lookAt(watchedShip.position);

            for (var i = 0; i < state['dust'].length; i++) {
                var p = state['dust'][i];

                var relX = p.position.x - state['camera'].position.x;
                var relY = p.position.y - state['camera'].position.y;

                if (relX < -1000) p.position.x += 2000;
                else if (relX > 1000) p.position.x -= 2000;

                if (relY < -1000) p.position.y += 2000;
                else if (relY > 1000) p.position.y -= 2000;

            }
        }


    });

    return SectorRenderer;

});