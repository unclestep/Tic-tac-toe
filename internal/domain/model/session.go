package model

import (
	"errors"
	"fmt"
)

type Session struct {
	UUID    string
	RulesID string
	Board   *Board
	Players []*Player
	Turn    int
	Winner  string
	State   State
	Params  *SessionParams
}

type State int

const (
	StateLobby State = iota
	StatePlaying
	StateGameOver
	StateUnknown
)

type SessionParams struct {
	Seed int64
}

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrPlayerExists   = errors.New("player already exists")
	ErrSessionFull    = errors.New("session is full")
	ErrMarkTaken      = errors.New("mark is taken")
)

func NewSession(sessionUUID string, params *SessionParams, rules *Rules) *Session {
	return &Session{
		UUID:    sessionUUID,
		RulesID: rules.UUID,
		Board:   NewBoard(rules.BoardWidth, rules.BoardHeight),
		Players: make([]*Player, 0, 2),
		State:   StateLobby,
		Params:  params,
	}
}

func (s *Session) Clone() *Session {
	return &Session{
		UUID:    s.UUID,
		RulesID: s.RulesID,
		Board:   s.Board.Clone(),
		Players: s.ClonePlayers(),
		Turn:    s.Turn,
		Winner:  s.Winner,
		State:   s.State,
		Params: &SessionParams{
			Seed: s.Params.Seed,
		},
	}
}

func (s *Session) ClonePlayers() []*Player {
	clone := make([]*Player, len(s.Players))
	for i, player := range s.Players {
		clone[i] = player.Clone()
	}
	return clone
}

func (s *Session) IsPlayerExist(playerUUID string) bool {
	for _, player := range s.Players {
		if player.UUID == playerUUID {
			return true
		}
	}
	return false
}

func (s *Session) IsFull() bool {
	return len(s.Players) == 2
}

func (s *Session) GetAvailableMarks() []Mark {
	if len(s.Players) == 0 {
		return []Mark{MarkX, MarkO}
	} else if len(s.Players) == 2 {
		return []Mark{}
	}
	return []Mark{s.Players[0].Mark.Opposite()}
}

func (s *Session) AddPlayer(player *Player) error {
	if s.IsFull() {
		return fmt.Errorf("add player: %w", ErrSessionFull)
	}
	for _, p := range s.Players {
		if p.UUID == player.UUID {
			return fmt.Errorf("add player: %w", ErrPlayerExists)
		}
		if p.Mark == player.Mark {
			return fmt.Errorf("add player: %w", ErrMarkTaken)
		}
	}
	s.Players = append(s.Players, player)
	return nil
}

func (s *Session) RemovePlayer(UUID string) error {
	for i, player := range s.Players {
		if player.UUID == UUID {
			if i == 0 {
				s.Players = s.Players[1:]
			} else {
				s.Players = s.Players[:1]
			}
			return nil
		}
	}
	return fmt.Errorf("remove player (session %s, player %s): %w", s.UUID, UUID, ErrPlayerNotFound)
}

func (s *Session) GetTurnPlayer() (*Player, error) {
	if len(s.Players) == 0 {
		return nil, fmt.Errorf("get turn player: %w", ErrPlayerNotFound)
	}
	return s.Players[s.Turn%len(s.Players)], nil
}
