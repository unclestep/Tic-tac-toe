package service

import (
	"tictactoe/internal/domain/model"
)

type Advancer struct {
	botLogic *BotLogic
}

func NewAdvancer(botLogic *BotLogic) *Advancer {
	return &Advancer{
		botLogic: botLogic,
	}
}

func (a *Advancer) Advance(session *model.Session) error {
	for session.State != model.StateGameOver {
		player := session.GetTurnPlayer()
		if player == nil || !player.IsBot {
			break
		}

		err := a.botLogic.BestMove(session, player.Mark)
		if err != nil {
			return err
		}
	}
	return nil
}
