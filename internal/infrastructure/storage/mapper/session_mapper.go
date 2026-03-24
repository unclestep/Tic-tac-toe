package mapper

import (
	domain "tictactoe/internal/domain/model"
	record "tictactoe/internal/infrastructure/storage/model"
)

func ToSessionRecord(s *domain.Session) record.SessionRecord {
	return record.SessionRecord{
		UUID:    s.UUID,
		RulesID: s.RulesID,
		Board:   toBoardRecord(s.Board),
		Players: convPlayersToRecords(s.Players),
		Turn:    s.Turn,
		Winner:  s.Winner,
		State:   stateToString(s.State),
		Seed:    s.Params.Seed,
	}
}

func toBoardRecord(b *domain.Board) record.BoardRecord {
	r := record.BoardRecord{
		Width:  b.Width,
		Height: b.Height,
		Cells:  make([]int8, b.Width*b.Height),
	}

	for i, cell := range b.CloneCells() {
		r.Cells[i] = int8(cell)
	}

	return r
}

func convPlayersToRecords(players []*domain.Player) []record.PlayerRecord {
	records := make([]record.PlayerRecord, len(players))
	for i, player := range players {
		records[i] = toPlayerRecord(player)
	}
	return records
}

func toPlayerRecord(p *domain.Player) record.PlayerRecord {
	return record.PlayerRecord{
		ID:   p.UUID,
		Name: p.Name,
		Mark: markToString(p.Mark),
	}
}

func markToString(mark domain.Mark) string {
	switch mark {
	case domain.X:
		return "X"
	case domain.O:
		return "O"
	default:
		return "Empty"
	}
}

func stateToString(state domain.State) string {
	switch state {
	case domain.StateLobby:
		return "Lobby"
	case domain.StatePlaying:
		return "Playing"
	case domain.StateGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

func ToSessionDomain(r record.SessionRecord) *domain.Session {
	return &domain.Session{
		UUID:    r.UUID,
		RulesID: r.RulesID,
		Board:   toBoardDomain(r.Board),
		Players: convPlayersToDomain(r.Players),
		Turn:    r.Turn,
		Winner:  r.Winner,
		State:   stringToState(r.State),
		Params:  &domain.SessionParams{Seed: r.Seed},
	}
}

func toBoardDomain(r record.BoardRecord) *domain.Board {
	cells := make([]domain.Mark, len(r.Cells))
	for i, cell := range r.Cells {
		cells[i] = domain.Mark(cell)
	}
	return domain.NewBoardFromCells(r.Width, r.Height, cells)
}

func convPlayersToDomain(records []record.PlayerRecord) []*domain.Player {
	players := make([]*domain.Player, len(records))
	for i, r := range records {
		players[i] = toPlayerDomain(r)
	}
	return players
}

func toPlayerDomain(r record.PlayerRecord) *domain.Player {
	return &domain.Player{
		UUID: r.ID,
		Name: r.Name,
		Mark: stringToMark(r.Mark),
	}
}

func stringToMark(s string) domain.Mark {
	switch s {
	case "X":
		return domain.X
	case "O":
		return domain.O
	default:
		return domain.Empty
	}
}

func stringToState(s string) domain.State {
	switch s {
	case "Lobby":
		return domain.StateLobby
	case "Playing":
		return domain.StatePlaying
	case "GameOver":
		return domain.StateGameOver
	default:
		return domain.StateLobby
	}
}
