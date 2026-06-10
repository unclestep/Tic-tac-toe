package dto

import (
	"tictactoe/pkg/geometry"
)

type CreateRequest struct {
	Rules Rules `json:"rules"`
	Seed  int64 `json:"seed"`
}

type Rules struct {
	BoardWidth  int `json:"board_width"`
	BoardHeight int `json:"board_height"`
	WinLength   int `json:"win_length"`
}

type ConnectRequest struct {
	SessionUUID string `json:"session_uuid"`
	PlayerName  string `json:"player_name"`
}

type StartRequest struct {
	SessionUUID string `json:"session_uuid"`
	PlayerUUID  string `json:"player_uuid"`
}

type MakeMoveRequest struct {
	SessionUUID string         `json:"session_uuid"`
	PlayerUUID  string         `json:"player_uuid"`
	MarkPos     geometry.Point `json:"mark_pos"`
}

type DisconnectRequest struct {
	SessionUUID string `json:"session_uuid"`
	PlayerUUID  string `json:"player_uuid"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
