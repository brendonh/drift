define([
], function() {

    var socket;
    if ("WebSocket" in window) {
        return window.WebSocket;
    } else if ("MozWebSocket" in window) {
        return window.MozWebSocket;
    }

    alert("No websocket support!");
    return;
});