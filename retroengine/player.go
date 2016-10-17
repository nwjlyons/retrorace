package retroengine

type Player struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Role  Role   `json:"role"`
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Role: Normal,
	}
}

func NewAdminPlayer(name string) *Player {
	return &Player{
		Name: name,
		Role: Admin,
	}
}

func (p *Player) IsAdmin() bool {
	return p.Role == Admin
}
