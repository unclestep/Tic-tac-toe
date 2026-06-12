package model

type Rules struct {
	BoardWidth  int
	BoardHeight int
	WinLength   int
}

func NewDefaultRules() *Rules {
	return &Rules{
		BoardWidth:  3,
		BoardHeight: 3,
		WinLength:   3,
	}
}

func (r *Rules) Clone() *Rules {
	if r == nil {
		return nil
	}
	return &Rules{
		BoardWidth:  r.BoardWidth,
		BoardHeight: r.BoardHeight,
		WinLength:   r.WinLength,
	}
}
