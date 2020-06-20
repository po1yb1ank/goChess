let socket = new WebSocket("wss://gochess.herokuapp.com/ws")
socket.onopen = () =>{
    var myCol
    var gs = false
    console.log("success connect from client")
    //socket.send("hi from client")
    socket.onclose = (event) =>{
        alert(event)
    }
    socket.onerror = (error) =>{
        alert(error)
    }
    socket.onmessage = (msg) =>{
        //alert(msg.data.toString())
        if(msg.data.toString() ==='gs'){
            alert("Both connected.Game started")
            gs = true
        }
        if(msg.data.toString() === 'w' || msg.data.toString() === 'b'){
            myCol = msg.data.toString()
                alert("Connected.")
                alert("your color is "+myCol)
            //game.setTurn(myCol)
        }else{
            game.load(msg.data.toString())
            board.position(msg.data.toString())
            m = msg.data.toString().split(" ")
            game.setTurn(m[1])
            updateStatus()
        }
    }

    var board = null
    var game = new Chess()
    var $status = $('#status')
    var $fen = $('#fen')
    var $pgn = $('#pgn')
    var m
    function onDragStart (source, piece, position, orientation) {
        // do not pick up pieces if the game is over
        if (game.game_over()) return false

        // only pick up pieces for the side to move
        if ((game.turn() !== myCol)||(gs == false)) {
            return false
        }

    }

    function onDrop (source, target) {
        // see if the move is legal
        var move = game.move({
            from: source,
            to: target,
            promotion: 'q' // NOTE: always promote to a queen for example simplicity
        })
        // illegal move
        if (move === null) return 'snapback'
        //updateStatus()

    }

    // update the board position after the piece snap
    // for castling, en passant, pawn promotion
    function onSnapEnd () {
        //board.position(game.fen())
        socket.send(game.fen())
        //alert(game.turn())
    }

    function updateStatus () {
        var status = ''
        var moveColor = 'White'
        if (game.turn() === 'b') {
            moveColor = 'Black'
        }

        // checkmate?
        if (game.in_checkmate()) {
            status = 'Game over, ' + moveColor + ' is in checkmate.'
        }

        // draw?
        else if (game.in_draw()) {
            status = 'Game over, drawn position'
        }

        // game still on
        else {
            status = moveColor + ' to move'

            // check?
            if (game.in_check()) {
                status += ', ' + moveColor + ' is in check'
            }
        }

        $status.html(status)
        $fen.html(game.fen())
        $pgn.html(game.pgn())
    }

    var config = {
        pieceTheme: "/static/styles/img/chesspieces/wikipedia/{piece}.png",
        draggable: true,
        position: 'start',
        onDragStart: onDragStart,
        onDrop: onDrop,
        onSnapEnd: onSnapEnd
    }
    board = Chessboard('myBoard', config)

    //updateStatus()
}