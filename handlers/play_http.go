package handlers

import (
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
)

func playHTTP(ctx *macaron.Context, game *retroengine.Game, player *retroengine.Player) {
	ctx.Data["game"] = game
	// If player is nil then the user is a spectator
	ctx.Data["player"] = player
	ctx.Data["winningCount"] = retroengine.WinningCount
	ctx.HTML(200, "play")
}
