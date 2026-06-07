package main

import (
	"go.uber.org/fx"
	_ "tictactoe/docs"
	"tictactoe/internal/di"
)

// @title           Tic-Tac-Toe API
// @version         1.0
// @description     Tic-Tac-Toe Game API
// @host            localhost:12121
// @BasePath        /
func main() {
	fx.New(di.TicTacToe).Run()
}
