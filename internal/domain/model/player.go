package model

type Player struct {
	UUID     string
	UserUUID string
	Name     string
	Mark     Mark
}

type Mark int8

const (
	MarkEmpty Mark = 0
	MarkX     Mark = 1
	MarkO     Mark = -1
)

func (m Mark) Opposite() Mark {
	switch m {
	case MarkX:
		return MarkO
	case MarkO:
		return MarkX
	}
	return MarkO
}

func (m Mark) String() string {
	switch m {
	case MarkEmpty:
		return " "
	case MarkX:
		return "X"
	case MarkO:
		return "O"
	default:
		return "?"
	}
}

func NewPlayer(UUID, UserUUID, name string, mark Mark) *Player {
	return &Player{
		UUID:     UUID,
		UserUUID: UserUUID,
		Name:     name,
		Mark:     mark,
	}
}

func NewBot(mark Mark) *Player {
	return &Player{
		UUID:     "BOT",
		UserUUID: "BOT",
		Name:     "BOT",
		Mark:     mark,
	}
}

func (p *Player) Clone() *Player {
	if p == nil {
		return nil
	}
	return &Player{
		UUID:     p.UUID,
		UserUUID: p.UserUUID,
		Name:     p.Name,
		Mark:     p.Mark,
	}
}
