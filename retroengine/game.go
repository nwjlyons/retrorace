package retroengine

import (
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

const (
	minNumPlayers = 2
	WinningCount  = 50
)

var (
	ticks = [...]string{"5", "4", "3", "2", "1"}
)

type Game struct {
	sync.Mutex
	Key        string                       `json:"key"`
	WebSockets map[*websocket.Conn]struct{} `json:"-"`
	Players    []*Player                    `json:"players"`
	State      State                        `json:"state"`
}

func NewGame(key string, admin string) *Game {
	game := Game{
		State:      WaitingForPlayers,
		Key:        key,
		WebSockets: make(map[*websocket.Conn]struct{}, 0),
	}
	game.Players = append(game.Players, NewAdminPlayer(admin))
	return &game
}

func (g *Game) AddSocket(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()
	g.WebSockets[conn] = struct{}{}
}

func (g *Game) RemoveSocket(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()
	delete(g.WebSockets, conn)
}

func (g *Game) BroadcastState() {
	g.Lock()
	defer g.Unlock()
	g.broadcastState()
}

func (g *Game) broadcastState() {
	for conn, _ := range g.WebSockets {
		// Sending self
		if err := conn.WriteJSON(GameStateWebSocketResponse{MsgType: "state", Game: g}); err != nil {
			// If failing to write on socket. Remove it.
			delete(g.WebSockets, conn)
		}
	}
}

func (g *Game) broadcast(v interface{}) {
	for conn, _ := range g.WebSockets {
		// Sending self
		if err := conn.WriteJSON(v); err != nil {
			// If failing to write on socket. Remove it.
			delete(g.WebSockets, conn)
		}
	}
}

func (g *Game) GetPlayer(name string) *Player {
	g.Lock()
	defer g.Unlock()
	return g.getPlayer(name)
}

func (g *Game) getPlayer(name string) *Player {
	for _, p := range g.Players {
		if strings.EqualFold(p.Name, name) {
			return p
		}
	}
	return nil
}

func (g *Game) AddPlayer(name string) *Player {
	g.Lock()
	defer g.Unlock()
	if p := g.getPlayer(name); p != nil {
		return nil
	} else {
		p := NewPlayer(name)
		g.Players = append(g.Players, p)
		g.broadcastState()
		return p
	}
}

func (g *Game) CloseToNewPlayers(player *Player) {
	g.Lock()
	defer g.Unlock()
	if g.State == WaitingForPlayers && player.IsAdmin() && len(g.Players) >= minNumPlayers {
		g.State = WaitingForCountdown
		g.broadcastState()
	}
}

func (g *Game) Start(player *Player) {
	g.Lock()
	defer g.Unlock()
	if g.State == WaitingForCountdown && player.IsAdmin() {
		g.State = CountingDown
		g.broadcastState()

		// Five second countdown.
		for _, tick := range ticks {
			g.broadcast(CountDownWebSocketResponse{
				MsgType: "countdown",
				Tick:    tick,
			})
			time.Sleep(time.Second)
		}
		g.broadcast(CountDownWebSocketResponse{
			MsgType: "countdown",
			Tick:    "Go!",
		})
		// Show Go! for a third of a second.
		time.Sleep(time.Second / 3)
		g.State = Started
		g.broadcastState()
	}
}

func (g *Game) Increment(player *Player) {
	g.Lock()
	defer g.Unlock()
	if g.State == Started {
		player.Count++
		if player.Count >= WinningCount {
			g.State = Finished
		}
		g.broadcastState()
	}
}

func (g *Game) Reset(player *Player) {
	g.Lock()
	defer g.Unlock()
	if player.IsAdmin() && g.State == Finished {
		for _, p := range g.Players {
			p.Count = 0
		}
		g.State = WaitingForCountdown
		g.broadcastState()
	}
}
