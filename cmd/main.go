package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	_ "tictactoe/docs"
	"tictactoe/internal/di"
)

func main() {
	_ = godotenv.Load(".env")
	fx.New(di.TicTacToe).Run()
}
