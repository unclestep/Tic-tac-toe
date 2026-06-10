package mapper

import (
	"fmt"
	"strings"

	"tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
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
			m, err := s.Board.GetMark(p)
			if err != nil {
				panic(fmt.Sprintf("to session response: %s", err))
			}
			sb.WriteByte(markToByte(m))
		}
		board[i] = sb.String()
	}

	turnPlayer, _ := s.GetTurnPlayer() //nolint:errcheck

	resp := dto.SessionResponse{
		SessionUUID: s.UUID,
		State:       stateToString(s.State),
		Board:       board,
		Winner:      playerToResponse(s.Winner),
		TurnPlayer:  playerToResponse(turnPlayer),
		Players:     playersToResponse(s.Players),
	}

	return resp
}

func ToUserReposnse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		UUID: user.UUID,
	}
}

func playersToResponse(players []*model.Player) []*dto.Player {
	ps := make([]*dto.Player, 0, len(players))
	for _, player := range players {
		ps = append(ps, playerToResponse(player))
	}
	return ps
}

func playerToResponse(player *model.Player) *dto.Player {
	if player == nil {
		return nil
	}
	return &dto.Player{
		UUID: player.UUID,
		Name: player.Name,
		Mark: string(markToByte(player.Mark)),
	}
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
	case model.MarkEmpty:
		return ' '
	case model.MarkX:
		return 'X'
	case model.MarkO:
		return 'O'
	default:
		return '?'
	}
}

func ToCreateCommand(req *dto.CreateRequest) *port.CreateCommand {
	return &port.CreateCommand{
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
		SessionUUID: req.SessionUUID,
		PlayerName:  req.PlayerName,
	}
}

func ToStartCommand(req *dto.StartRequest) *port.StartCommand {
	return &port.StartCommand{
		SessionUUID: req.SessionUUID,
		PlayerUUID:  req.PlayerUUID,
	}
}

func ToMakeMoveCommand(req *dto.MakeMoveRequest) *port.MakeMoveCommand {
	return &port.MakeMoveCommand{
		SessionUUID: req.SessionUUID,
		PlayerUUID:  req.PlayerUUID,
		MarkPos:     req.MarkPos,
	}
}

func ToDisconnectCommand(req *dto.DisconnectRequest) *port.DisconnectCommand {
	return &port.DisconnectCommand{
		SessionUUID: req.SessionUUID,
		PlayerUUID:  req.PlayerUUID,
	}
}

func ToSignUpCommand(req *dto.SignUpRequest) *usecase.SignUpCommand {
	return &usecase.SignUpCommand{
		Login:    req.Login,
		Password: req.Password,
	}
}

func ToSignInCommand(req *dto.SignInRequest) *usecase.SignInCommand {
	return &usecase.SignInCommand{
		Login:    req.Login,
		Password: req.Password,
	}
}
