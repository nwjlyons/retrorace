package retroengine

import "sync"

type Games struct {
	sync.Mutex
	Games []*Game
}

var GameStore = &Games{}

func (gs *Games) Get(key string) *Game {
	for _, g := range gs.Games {
		if g.Key == key {
			return g
		}
	}
	return nil
}

func (gs *Games) New(admin string) *Game {
	gs.Lock()
	defer gs.Unlock()
	for {
		key := NewKey()
		if gs.Get(key) == nil {
			g := NewGame(key, admin)
			gs.Games = append(gs.Games, g)
			return g
		}
	}
}
