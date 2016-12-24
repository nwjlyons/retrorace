package handlers

import (
	"strings"

	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	if ctx.Req.Method == "POST" {
		playerName := strings.TrimSpace(ctx.Req.FormValue("playerName"))
		if playerName == "" {
			ctx.Data["err"] = "This field is required."
			ctx.HTML(400, "index")
		} else if len(playerName) > 5 {
			ctx.Data["err"] = "Name too long. Five characters of less."
			ctx.HTML(400, "index")
		} else {
			game := retroengine.GameStore.New(playerName)
			gamePath := ctx.URLFor("play", "key", game.Key)
			ctx.SetCookie("playerName", playerName, 0, gamePath)
			ctx.Redirect(gamePath)
		}
	} else {
		ctx.HTML(200, "index")
	}
}
