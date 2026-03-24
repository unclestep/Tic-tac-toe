package main

import (
	"go.uber.org/fx"
	"tictactoe/internal/di"
)

func main() {
	fx.New(di.Module).Run()
}
