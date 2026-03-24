package mapper

import (
	"strings"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

func ToSessionResponse(s *model.Session) dto.SessionResponse {
	board := make([]string, s.Board.Height)
	var sb strings.Builder
	for i := range s.Board.Height {
		sb.Reset()
		for j := range s.Board.Width {
			p := geometry.NewPoint(j, i)
			sb.WriteByte(markToByte(s.Board.GetMark(p)))
		}
		board[i] = sb.String()
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

func markToByte(mark model.Mark) byte {
	switch mark {
	case model.Empty:
		return ' '
	case model.X:
		return 'X'
	case model.O:
		return 'O'
	default:
		return '?'
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
