package handlers

import (
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
	"strings"
)

func join(ctx *macaron.Context, game *retroengine.Game) {
	ctx.Data["spectatePath"] = ctx.URLFor("spectate", "key", game.Key)

	if game.State == retroengine.WaitingForPlayers && len(game.Players) <= 5 {
		ctx.Data["canStillJoin"] = true
	} else {
		ctx.Data["canStillJoin"] = false
	}

	if ctx.Data["canStillJoin"].(bool) && ctx.Req.Method == "POST" {
		playerName := strings.TrimSpace(ctx.Req.FormValue("playerName"))
		if playerName == "" {
			ctx.Data["err"] = "This field is required."
			ctx.HTML(400, "join")
		} else if len(playerName) > 5 {
			ctx.Data["err"] = "Name too long. Five characters of less."
			ctx.HTML(400, "index")
		} else {
			if p := game.AddPlayer(playerName); p == nil {
				ctx.Data["err"] = "Player with this name already exists."
				ctx.HTML(400, "join")
			} else {
				gamePath := ctx.URLFor("play", "key", game.Key)
				ctx.SetCookie("playerName", playerName, 0, gamePath)
				ctx.Redirect(gamePath)
			}
		}
	} else {
		ctx.HTML(200, "join")
	}
}
