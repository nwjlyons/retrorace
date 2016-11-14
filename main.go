package main

import (
	"flag"
	"github.com/go-macaron/pongo2"
	"github.com/nwjlyons/retrorace/handlers"
	"gopkg.in/macaron.v1"
)

var (
	port = flag.Int("p", 8000, "port")
)

func main() {
	flag.Parse()
	m := macaron.New()

	// Middleware
	m.Use(pongo2.Pongoer(pongo2.Options{
		IndentJSON: true,
	}))
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static", macaron.StaticOptions{
		Prefix:      "static",
		SkipLogging: true,
	}))

	// Routes
	m.Get("/_stats", handlers.Stats).Name("stats")
	m.Get("/:key/spectate", handlers.Spectate).Name("spectate")
	m.Combo("/:key").Get(handlers.Play).Post(handlers.Play).Name("play")
	m.Combo("/").Get(handlers.Index).Post(handlers.Index).Name("index")

	// Lets Race!
	m.Run("localhost", *port)
}
