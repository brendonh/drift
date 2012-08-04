var state;

function setup() {
    window.requestAnimFrame = (function(){
        return  window.requestAnimationFrame   || 
            window.webkitRequestAnimationFrame || 
            window.mozRequestAnimationFrame    || 
            window.oRequestAnimationFrame      || 
            window.msRequestAnimationFrame     || 
            function( callback ){
                window.setTimeout(callback, 1000 / 60);
            };
    })();

    var canvas = $('#box').get(0);

    var scene = new THREE.Scene();

    var camera = new THREE.PerspectiveCamera( 60, window.innerWidth / window.innerHeight, 1, 10000 );
    camera.position.z = 1000;
    scene.add( camera );

    // var geometry = new THREE.CubeGeometry( 
    //     200, 200, 200,
    //     undefined, undefined, undefined,
    //     undefined,
    //     {nz: false}
    // );

    var geometry = new CubeSet().toGeometry();

    var material = new THREE.MeshBasicMaterial( { color: 0x999999 } );
    
    var mesh = new THREE.Mesh( geometry, material );
    scene.add( mesh );
    
    var renderer = new THREE.CanvasRenderer({'canvas': canvas});
    renderer.setSize( 800, 600 );
    
    state = {
        'frames': 0,
        'lastTime': new Date().getTime(),

        'canvas': canvas,
        'scene': scene,
        'camera': camera,        
        'mesh': mesh,
        'renderer': renderer
    };

}


function go() {
    setup();
    window.requestAnimFrame(frame);
}


function frame() {
    state['frames'] += 1;

    // state['mesh'].rotation.x = 0.5;
    // state['mesh'].rotation.y = 0.5;
    state['mesh'].rotation.x += 0.01;
    state['mesh'].rotation.y += 0.02;

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