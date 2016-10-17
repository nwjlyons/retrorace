package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
)

func Play(ctx *macaron.Context) {
	game := retroengine.GameStore.Get(ctx.Params("key"))
	if game == nil {
		ctx.Redirect("/")
		return
	}

	playerName := ctx.GetCookie("playerName")
	if playerName == "" {
		join(ctx, game)
		return
	}

	player := game.GetPlayer(playerName)
	if player == nil {
		ctx.Status(403)
		return
	}

	// route to a handler based on request type.
	if websocket.IsWebSocketUpgrade(ctx.Req.Request) {
		// websocket request
		playWS(ctx, game, player)
	} else {
		// http request
		playHTTP(ctx, game, player)
	}
}
