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
	Bots    int
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
	ErrPlayerNotFound      = errors.New("player not found")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrCantCreateSession   = errors.New("cant create session")
	ErrSessionFull         = errors.New("session is full")
	ErrGameAlreadyOver     = errors.New("game is already over")
)

func NewSession(sessionUUID string, params *SessionParams, rules *Rules) (*Session, error) {
	if sessionUUID == "" {
		return nil, fmt.Errorf("%w: session id is invalid", ErrCantCreateSession)
	}
	if params == nil {
		return nil, fmt.Errorf("%w: params is nil", ErrCantCreateSession)
	}
	if rules == nil {
		return nil, fmt.Errorf("%w: rules is nil", ErrCantCreateSession)
	}
	if rules.UUID == "" {
		return nil, fmt.Errorf("%w: rules have invalid id", ErrCantCreateSession)
	}

	board, err := NewBoard(rules.BoardWidth, rules.BoardHeight)
	if err != nil {
		return nil, err
	}

	return &Session{
		UUID:    sessionUUID,
		RulesID: rules.UUID,
		Board:   board,
		Players: make([]*Player, 0, 2),
		State:   StateLobby,
		Params:  params,
	}, nil
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
		if player.UUID == playerUUID && !player.IsBot {
			return true
		}
	}
	return false
}

const (
	BotUUID = "BotUUID"
	BotName = "Bot"
)

func (s *Session) AddBot() {
	bot := NewPlayer(BotUUID, BotName)
	bot.IsBot = true
	s.Bots++
	s.Players = append(s.Players, bot)
}

func (s *Session) AddPlayer(playerUUID, playerName string) error {
	var bot *Player
	for _, player := range s.Players {
		if player.UUID == playerUUID && player.IsBot {
			player.IsBot = false
			s.Bots--
			return nil
		} else if player.UUID == playerUUID {
			return fmt.Errorf("%w: session %s, player %s", ErrPlayerAlreadyExists, s.UUID, playerUUID)
		} else if player.IsBot {
			bot = player
		}
	}

	if len(s.Players)-s.Bots >= 2 {
		return fmt.Errorf("%w: session %s", ErrSessionFull, s.UUID)
	}

	if bot != nil {
		bot.UUID = playerUUID
		bot.Name = playerName
		bot.IsBot = false
		s.Bots--
		return nil
	}

	player := NewPlayer(playerUUID, playerName)
	if len(s.Players) == 0 {
		player.Mark = X
	} else {
		player.Mark = O
	}
	s.Players = append(s.Players, player)

	return nil
}

func (s *Session) HidePlayer(id string) error {
	for _, player := range s.Players {
		if player.UUID == id {
			if player.IsBot {
				break
			}
			player.IsBot = true
			s.Bots++
			return nil
		}
	}

	return fmt.Errorf("%w: player %v", ErrPlayerNotFound, id)
}

func (s *Session) RemovePlayer(playerUUID string) error {
	removeInd := -1
	for i, player := range s.Players {
		if player.UUID == playerUUID {
			removeInd = i
			break
		}
	}
	if removeInd == -1 {
		return fmt.Errorf("%w: session %s, player %s", ErrPlayerNotFound, s.UUID, playerUUID)
	}

	last := len(s.Players) - 1
	s.Players[removeInd] = s.Players[last]
	s.Players = s.Players[:last]
	return nil
}

func (s *Session) GetTurnPlayer() *Player {
	if len(s.Players) == 0 {
		return nil
	}
	return s.Players[s.Turn%len(s.Players)]
}

func (s *Session) DetermineWinner(mark Mark) string {
	if mark == Empty {
		return "Draw"
	}

	for _, player := range s.Players {
		if player.Mark == mark {
			return player.Name
		}
	}

	return "Bot"
}
