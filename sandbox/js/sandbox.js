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

    state = {
        'frames': 0,
        'lastTime': new Date().getTime(),

        'canvas': $('#box').get(0),

        'quads': new Array()
    };

    state['ctx'] = state['canvas'].getContext('2d');


}


function createQuads(count) {
    var rand255 = function() { return Math.round(Math.random() * 255); }
    for (var i = 0; i < count; i++) {
        var fillStyle = 'rgb(' + rand255() + ', ' + rand255() + ', ' + rand255() + ')';
        state['quads'].push([Math.random() * 800, Math.random() * 600, 0, fillStyle]);
    }
}

function moveQuads() {
    var clampVary = function(val, vary, min, max) { 
        val += (Math.random() * vary) - (vary / 2);
        return Math.max(min, Math.min(max, val)); 
    }
    for (var i = 0; i < state['quads'].length; i++) {
        var quad = state['quads'][i];
        quad[0] = clampVary(quad[0], 2, 0, 800);
        quad[1] = clampVary(quad[1], 2, 0, 600);
        quad[2] = clampVary(quad[2], 0.4, 0, 6.28);
    }
}


function go() {
    setup();
    createQuads(1000);
    window.requestAnimFrame(frame);
}


function frame() {
    state['frames'] += 1;

    info();

    moveQuads();

    var ctx = state['ctx'];

    ctx.fillStyle = 'rgb(0, 0, 0)';
    ctx.fillRect(0, 0, 800, 600);

    ctx.save();

    for (var i = 0; i < state['quads'].length; i++) {
        var quad = state['quads'][i];

        ctx.fillStyle = quad[3];
        
        ctx.translate(quad[0], quad[1]);
        ctx.rotate(quad[2]);
        
        ctx.beginPath();
        ctx.moveTo(-5, -5);
        ctx.lineTo( 5, -5);
        ctx.lineTo( 5,  5);
        ctx.lineTo(-5,  5);
        ctx.closePath();
        ctx.fill();

        ctx.restore();
        ctx.save();

    }

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