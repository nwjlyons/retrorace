package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func playWS(ctx *macaron.Context, game *retroengine.Game, player *retroengine.Player) {

	// Upgrade the request to a websocket connection
	conn, err := upgrader.Upgrade(ctx.Resp, ctx.Req.Request, nil)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	game.AddSocket(conn)
	game.BroadcastState()

	// Loop forever
	for {
		// ReadMessage is blocking. It will pause the execution of this for loop until a message is
		// received.
		messageType, action, err := conn.ReadMessage()

		if err != nil {
			// For any type of error, remove the web socket from the game and break out of the loop.
			game.RemoveSocket(conn)
			return
		}

		// Only interested in text messages.
		// nil players are spectators. They cannot perform any actions, so ignore any message they send.
		if player != nil && messageType == websocket.TextMessage {
			switch string(action) {
			case CloseToNewPlayers.String():
				game.CloseToNewPlayers(player)
			case Start.String():
				game.Start(player)
			case Increment.String():
				game.Increment(player)
			case Reset.String():
				game.Reset(player)
			}
		}
	}
}
