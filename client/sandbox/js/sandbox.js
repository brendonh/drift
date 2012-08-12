var state;

function setup() {

    var $container = $('#box');

    var scene = new THREE.Scene();

    var camera = new THREE.PerspectiveCamera( 60, 800 / 600, 1, 2000 );
    camera.position.z = 1000;
    scene.add( camera );

    var renderer;
    if (window.location.search.indexOf("webgl") != -1) {
        renderer = new THREE.WebGLRenderer();
        $("#info .renderer").html("WebGL");
    } else {
        renderer = new THREE.CanvasRenderer();
        $("#info .renderer").html("Canvas");
    }
    renderer.setSize( 800, 600 );
    $(renderer.domElement).css("height", "600px").css("width", "800px").css("background-color", "black");
    $container.empty().append(renderer.domElement);

    state = {
        'frames': 0,
        'lastTime': new Date().getTime(),

        'width': 800,
        'height': 600,

        'scene': scene,
        'camera': camera,        
        'renderer': renderer,

        'rotation': 0,
        'flip': 0,
        'moving': false,
        'rotateLeft': false,
        'rotateRight': false
    };

}

function addDust() {

    var projector = new THREE.Projector();
    var widthHalf = state['width'] / 2;
    var heightHalf = state['height'] / 2;

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

    for(var p = 0; p < count; p++) {

        var particle = new THREE.Particle(pMaterial);

        particle.position.x = -1000 + Math.random() * 2000;
        particle.position.y = -1000 + Math.random() * 2000;
        particle.position.z = Math.random() * 500;

        particles.push(particle);
        state['scene'].add(particle);
    }

    state['dust'] = particles;
}

function addShip() {

    var colors = {
        'X': 0x993333,
        'x': 0x999999,
        'o': 0x444444
    };

    var ship = new CubeSet(10);
    ship.addAscii([
        ' o o ',
        ' o o ',
        ' oxo ',
        ' xxx ',
        'xxXxx',
        'xxxxx',
        ' xox '
    ], colors);

    var mesh = ship.toMesh();
    mesh.useQuaternion = true;
    state['scene'].add( mesh );

    state['ship'] = mesh;
}

function go() {
    setup();
    addDust();
    addShip();

    var directionalLight = new THREE.DirectionalLight( 0xffffff, 1.0 );
	directionalLight.position.set( 0.5, 0.2, -1 );
    state['scene'].add(directionalLight);

    $("#info .faces").html(state['ship'].geometry.faces.length);

    window.onkeydown = function(e) {
        if (e.keyCode == 38) state['moving'] = true;
        else if (e.keyCode == 39) state['rotateRight'] = true;
        else if (e.keyCode == 37) state['rotateLeft'] = true;
    };

    window.onkeyup = function(e) {
        if (e.keyCode == 38) state['moving'] = false;
        else if (e.keyCode == 39) state['rotateRight'] = false;
        else if (e.keyCode == 37) state['rotateLeft'] = false;
    };

    window.requestAnimFrame(frame);
}


function frame() {
    state['frames'] += 1;

    state['camera'].position = state['ship'].position.clone();
    state['camera'].position.z = -700;
    state['camera'].lookAt(state['ship'].position);

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
        state['ship'].position.x -= speed * Math.sin(state['rotation']);
        state['ship'].position.y += speed * Math.cos(state['rotation']);
    }

    state['ship'].quaternion.setFromAxisAngle(new THREE.Vector3(0, 0, 1), state['rotation']);
    var flip = new THREE.Quaternion();
    flip.setFromAxisAngle(new THREE.Vector3(0, 1, 0), state['flip']);
    state['ship'].quaternion.multiplySelf(flip);

    for (var i = 0; i < state['dust'].length; i++) {
        var p = state['dust'][i];

        var relX = p.position.x - state['camera'].position.x;
        var relY = p.position.y - state['camera'].position.y;

        if (relX < -1000) p.position.x += 2000;
        else if (relX > 1000) p.position.x -= 2000;

        if (relY < -1000) p.position.y += 2000;
        else if (relY > 1000) p.position.y -= 2000;

    }

    state['renderer'].render( 
        state['scene'], 
        state['camera'] );

    info();

    window.requestAnimFrame(frame);
}


function info() {
    if (state['frames'] > 60) {
        var time = new Date().getTime();
        var delta = (time - state['lastTime']) / 1000;
        
        $("#info .fps").html( Math.round((state['frames'] / delta)) );
        
        state['frames'] = 0;
        state['lastTime'] = time;
    }
}