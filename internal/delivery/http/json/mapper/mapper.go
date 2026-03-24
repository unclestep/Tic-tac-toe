package mapper

import (
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/domain/model"
)

func ToSessionResponse(s *model.Session) dto.SessionResponse {
	cells := s.Board.CloneCells()
	board := make([]int8, len(cells))
	for i, cell := range cells {
		board[i] = int8(cell)
	}

	players := make([]string, 0, len(s.Players))
	for _, player := range s.Players {
		players = append(players, player.Name)
	}

	resp := dto.SessionResponse{
		SessionID: s.UUID,
		State:     stateToString(s.State),
		Board:     board,
		Players:   players,
		Winner:    s.Winner,
	}

	return resp
}

func stateToString(state model.State) string {
	switch state {
	case model.StateLobby:
		return "Lobby"
	case model.StatePlaying:
		return "Playing"
	case model.StateGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

func ToCreateCommand(req *dto.CreateRequest) *port.CreateCommand {
	return &port.CreateCommand{
		PlayerID:   req.PlayerID,
		PlayerName: req.PlayerName,
		Rules: &model.Rules{
			BoardWidth:  req.Rules.BoardWidth,
			BoardHeight: req.Rules.BoardHeight,
			WinLength:   req.Rules.WinLength,
		},
		SessionParams: &model.SessionParams{
			Seed: req.Seed,
		},
	}
}

func ToConnectCommand(req *dto.ConnectRequest) *port.ConnectCommand {
	return &port.ConnectCommand{
		SessionID:  req.SessionID,
		PlayerID:   req.PlayerID,
		PlayerName: req.PlayerName,
	}
}

func ToStartCommand(req *dto.StartRequest) *port.StartCommand {
	return &port.StartCommand{
		SessionID: req.SessionID,
		PlayerID:  req.PlayerID,
	}
}

func ToMakeMoveCommand(req *dto.MakeMoveRequest) *port.MakeMoveCommand {
	return &port.MakeMoveCommand{
		SessionID: req.SessionID,
		PlayerID:  req.PlayerID,
		MarkPos:   req.MarkPos,
	}
}

func ToDisconnectCommand(req *dto.DisconnectRequest) *port.DisconnectCommand {
	return &port.DisconnectCommand{
		SessionID: req.SessionID,
		PlayerID:  req.PlayerID,
	}
}
