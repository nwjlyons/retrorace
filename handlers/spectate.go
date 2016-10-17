package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
)

func Spectate(ctx *macaron.Context) {
	game := retroengine.GameStore.Get(ctx.Params("key"))
	if game == nil {
		ctx.Redirect("/")
		return
	}

	// Pass in nil for the player
	if websocket.IsWebSocketUpgrade(ctx.Req.Request) {
		playWS(ctx, game, nil)
	} else {
		playHTTP(ctx, game, nil)
	}
}
