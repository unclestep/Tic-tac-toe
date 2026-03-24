package dto

import (
	"tictactoe/pkg/geometry"
)

type CreateRequest struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Rules      Rules  `json:"rules"`
	Seed       int64  `json:"seed"`
}

type Rules struct {
	BoardWidth  int `json:"board_width"`
	BoardHeight int `json:"board_height"`
	WinLength   int `json:"win_length"`
}

type ConnectRequest struct {
	SessionID  string `json:"session_id"`
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

type StartRequest struct {
	SessionID string `json:"session_id"`
	PlayerID  string `json:"player_id"`
}

type MakeMoveRequest struct {
	SessionID string         `json:"session_id"`
	PlayerID  string         `json:"player_id"`
	MarkPos   geometry.Point `json:"mark_pos"`
}

type DisconnectRequest struct {
	SessionID string `json:"session_id"`
	PlayerID  string `json:"player_id"`
}
