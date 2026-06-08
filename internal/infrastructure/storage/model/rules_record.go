package model

type RulesRecord struct {
	UUID        string
	BoardWidth  int
	BoardHeight int
	WinLength   int
}

func (rd *RulesRecord) Clone() *RulesRecord {
	return &RulesRecord{
		UUID:        rd.UUID,
		BoardWidth:  rd.BoardWidth,
		BoardHeight: rd.BoardHeight,
		WinLength:   rd.WinLength,
	}
}
