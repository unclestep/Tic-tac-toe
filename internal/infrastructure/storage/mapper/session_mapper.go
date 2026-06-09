package mapper

import (
	"fmt"
	dmodel "tictactoe/internal/domain/model"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

func ToSessionStorage(s *dmodel.Session) *dsmodel.SessionRecord {
	return &dsmodel.SessionRecord{
		UUID:      s.UUID,
		RulesUUID: s.Rules.UUID,
		Board:     toBoardStorage(s.Board),
		Players:   convPlayersToStorage(s.Players),
		Turn:      s.Turn,
		Winner:    s.Winner,
		State:     stateToString(s.State),
		Seed:      s.Params.Seed,
	}
}

func toBoardStorage(b *dmodel.Board) *dsmodel.BoardRecord {
	r := dsmodel.BoardRecord{
		Width:  b.Width,
		Height: b.Height,
		Cells:  make([]int8, b.Width*b.Height),
	}

	for i, cell := range b.CloneCells() {
		r.Cells[i] = int8(cell)
	}

	return &r
}

func convPlayersToStorage(players []*dmodel.Player) []*dsmodel.PlayerRecord {
	records := make([]*dsmodel.PlayerRecord, len(players))
	for i, player := range players {
		records[i] = toPlayerStorage(player)
	}
	return records
}

func toPlayerStorage(p *dmodel.Player) *dsmodel.PlayerRecord {
	return &dsmodel.PlayerRecord{
		UUID: p.UUID,
		Name: p.Name,
		Mark: markToString(p.Mark),
	}
}

func markToString(mark dmodel.Mark) string {
	switch mark {
	case dmodel.MarkX:
		return "X"
	case dmodel.MarkO:
		return "O"
	case dmodel.MarkEmpty:
		return "Empty"
	default:
		return "Unknown"
	}
}

func stateToString(state dmodel.State) string {
	switch state {
	case dmodel.StateLobby:
		return "Lobby"
	case dmodel.StatePlaying:
		return "Playing"
	case dmodel.StateGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

func ToSessionDomain(r *dsmodel.SessionRecord) (*dmodel.Session, error) {
	players, err := convPlayersToDomain(r.Players)
	if err != nil {
		return nil, fmt.Errorf("to session domain: %w", err)
	}

	state, err := stringToState(r.State)
	if err != nil {
		return nil, fmt.Errorf("to session domain: %w", err)
	}

	return &dmodel.Session{
		UUID:    r.UUID,
		Board:   toBoardDomain(r.Board),
		Players: players,
		Turn:    r.Turn,
		Winner:  r.Winner,
		State:   state,
		Params:  &dmodel.SessionParams{Seed: r.Seed},
	}, nil
}

func toBoardDomain(r *dsmodel.BoardRecord) *dmodel.Board {
	cells := make([]dmodel.Mark, len(r.Cells))
	for i, cell := range r.Cells {
		cells[i] = dmodel.Mark(cell)
	}
	return dmodel.NewBoardFromCells(r.Width, r.Height, cells)
}

func convPlayersToDomain(records []*dsmodel.PlayerRecord) ([]*dmodel.Player, error) {
	players := make([]*dmodel.Player, 0, len(records))
	for _, r := range records {
		player, err := toPlayerDomain(r)
		if err != nil {
			return nil, fmt.Errorf("conv players to domain: %w", err)
		}
		players = append(players, player)
	}
	return players, nil
}

func toPlayerDomain(r *dsmodel.PlayerRecord) (*dmodel.Player, error) {
	mark, err := stringToMark(r.Mark)
	if err != nil {
		return nil, fmt.Errorf("to player domain: %w", err)
	}

	return &dmodel.Player{
		UUID: r.UUID,
		Name: r.Name,
		Mark: mark,
	}, nil
}

func stringToMark(s string) (dmodel.Mark, error) {
	switch s {
	case "X":
		return dmodel.MarkX, nil
	case "O":
		return dmodel.MarkO, nil
	case "Empty":
		return dmodel.MarkEmpty, nil
	default:
		return 0, fmt.Errorf("unknown mark: %s", s)
	}
}

func stringToState(s string) (dmodel.State, error) {
	switch s {
	case "Lobby":
		return dmodel.StateLobby, nil
	case "Playing":
		return dmodel.StatePlaying, nil
	case "GameOver":
		return dmodel.StateGameOver, nil
	default:
		return dmodel.StateUnknown, fmt.Errorf("unknown state: %s", s)
	}
}
