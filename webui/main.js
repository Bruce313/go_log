window.onload = function() {
    var socket = io('http://localhost:3354/socket.io')
    socket.on('connect', function() {
        console.log('connect');
    });

    socket.on('new_log', function(data) {
        console.log('new_log', data);
        document.getElementById("logs").innnerHTML += 
            ("<li>" + data + "</li>");
    });

    socket.on('error', function(err) {
        console.error(err);
    });
}
