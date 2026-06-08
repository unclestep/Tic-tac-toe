package model

import "slices"

type SessionRecord struct {
	UUID      string
	RulesUUID string
	Board     *BoardRecord
	Players   []*PlayerRecord
	Bots      int
	Turn      int
	Winner    string
	State     string
	Seed      int64
}

func (sr *SessionRecord) Clone() *SessionRecord {
	var players []*PlayerRecord
	for _, p := range sr.Players {
		players = append(players, p.Clone())
	}

	return &SessionRecord{
		UUID:      sr.UUID,
		RulesUUID: sr.RulesUUID,
		Board:     sr.Board.Clone(),
		Players:   players,
		Bots:      sr.Bots,
		Turn:      sr.Turn,
		Winner:    sr.Winner,
		State:     sr.State,
		Seed:      sr.Seed,
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
	UUID string
	Name string
	Mark string
	Bot  bool
}

func (pr *PlayerRecord) Clone() *PlayerRecord {
	return &PlayerRecord{
		UUID: pr.UUID,
		Name: pr.Name,
		Mark: pr.Mark,
		Bot:  pr.Bot,
	}
}
