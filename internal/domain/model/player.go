package model

type Player struct {
	UUID string
	Name string
	Mark Mark
}

type Mark int8

const (
	MarkEmpty Mark = 0
	MarkX     Mark = 1
	MarkO     Mark = -1
)

func (m Mark) Opposite() Mark {
	if m == MarkX {
		return MarkO
	}
	return MarkX
}

func (m Mark) String() string {
	switch m {
	case MarkEmpty:
		return " "
	case MarkX:
		return "X"
	case MarkO:
		return "0"
	default:
		return "?"
	}
}

func NewPlayer(UUID, name string, mark Mark) *Player {
	return &Player{
		UUID: UUID,
		Name: name,
		Mark: mark,
	}
}

func (p *Player) Clone() *Player {
	return &Player{
		UUID: p.UUID,
		Name: p.Name,
		Mark: p.Mark,
	}
}
