package model

type SessionRecord struct {
	UUID    string
	RulesID string
	Board   BoardRecord
	Players []PlayerRecord
	Bots    int
	Turn    int
	Winner  string
	State   string
	Seed    int64
}

type BoardRecord struct {
	Width  int
	Height int
	Cells  []int8
}

type PlayerRecord struct {
	ID    string
	Name  string
	Mark  string
	IsBot bool
}
