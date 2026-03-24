package model

type Player struct {
	UUID string
	Name string
	Mark Mark
}

func NewPlayer(UUID string, name string) *Player {
	return &Player{
		UUID: UUID,
		Name: name,
	}
}

func (p *Player) Clone() *Player {
	return &Player{
		UUID: p.UUID,
		Name: p.Name,
		Mark: p.Mark,
	}
}
