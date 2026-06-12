package model

import "slices"

type SessionRecord struct {
	UUID    string
	Rules   *RulesRecord
	Board   *BoardRecord
	Players []*PlayerRecord
	Turn    int
	Winner  *PlayerRecord
	State   string
	Seed    int64
}

func (r *SessionRecord) Clone() *SessionRecord {
	var players []*PlayerRecord
	for _, p := range r.Players {
		players = append(players, p.Clone())
	}

	return &SessionRecord{
		UUID:    r.UUID,
		Rules:   r.Rules.Clone(),
		Board:   r.Board.Clone(),
		Players: players,
		Turn:    r.Turn,
		Winner:  r.Winner.Clone(),
		State:   r.State,
		Seed:    r.Seed,
	}
}

type RulesRecord struct {
	BoardWidth  int
	BoardHeight int
	WinLength   int
}

func (r *RulesRecord) Clone() *RulesRecord {
	if r == nil {
		return nil
	}
	return &RulesRecord{
		BoardWidth:  r.BoardWidth,
		BoardHeight: r.BoardHeight,
		WinLength:   r.WinLength,
	}
}

type BoardRecord struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []int8 `json:"cells"`
}

func (br *BoardRecord) Clone() *BoardRecord {
	return &BoardRecord{
		Width:  br.Width,
		Height: br.Height,
		Cells:  slices.Clone(br.Cells),
	}
}

type PlayerRecord struct {
	UUID     string
	UserUUID string
	Name     string
	Mark     string
}

func (pr *PlayerRecord) Clone() *PlayerRecord {
	return &PlayerRecord{
		UUID: pr.UUID,
		Name: pr.Name,
		Mark: pr.Mark,
	}
}
