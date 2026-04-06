package main

import (
	"go.uber.org/fx"
	_ "tictactoe/docs"
	"tictactoe/internal/di"
)

func main() {
	fx.New(di.Module).Run()
}
