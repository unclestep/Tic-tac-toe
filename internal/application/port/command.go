package port

import (
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

type ContextKey string

const UserUUIDKey ContextKey = "UserUUID"

type CreateCommand struct {
	Rules         *model.Rules
	SessionParams *model.SessionParams
}

type ConnectCommand struct {
	SessionUUID string
	PlayerName  string
}

type StartCommand struct {
	SessionUUID string
	PlayerUUID  string
}

type MakeMoveCommand struct {
	SessionUUID string
	PlayerUUID  string
	MarkPos     geometry.Point
}

type DisconnectCommand struct {
	SessionUUID string
	PlayerUUID  string
}
