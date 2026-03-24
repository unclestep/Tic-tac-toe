package port

import (
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type CreateCommand struct {
	PlayerID      string
	PlayerName    string
	Rules         *model.Rules
	SessionParams *model.SessionParams
}

type ConnectCommand struct {
	SessionID  string
	PlayerID   string
	PlayerName string
}

type StartCommand struct {
	SessionID string
	PlayerID  string
}

type MakeMoveCommand struct {
	SessionID string
	PlayerID  string
	MarkPos   geometry.Point
}

type DisconnectCommand struct {
	SessionID string
	PlayerID  string
}
