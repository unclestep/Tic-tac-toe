package model

type Rules struct {
	UUID                    string
	BoardWidth, BoardHeight int
	WinLength               int
}

func NewDefaultRules() *Rules {
	return &Rules{
		BoardWidth:  3,
		BoardHeight: 3,
		WinLength:   3,
	}
}

func (r *Rules) Clone() *Rules {
	return &Rules{
		UUID:        r.UUID,
		BoardWidth:  r.BoardWidth,
		BoardHeight: r.BoardHeight,
		WinLength:   r.WinLength,
	}
}
