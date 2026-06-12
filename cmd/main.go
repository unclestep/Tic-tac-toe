package main

import (
	"context"
	"log"
	"os"

	"tictactoe/docs"
	"tictactoe/internal/di"
	"tictactoe/internal/infrastructure/storage/postgres"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// @title           TicTacToe API
// @version         1.0
// @host            localhost:12121
// @BasePath        /
// @SecurityDefinitions.basic BasicAuth
func main() {
	_ = godotenv.Load(".env")
	docs.SwaggerInfo.Host = os.Getenv("TICTACTOE_HOST")
	if err := postgres.Migrate(context.Background(), os.Getenv("POSTGRES_DSN")); err != nil {
		log.Fatalf("migration failed: %s", err)
	}
	fx.New(di.TicTacToe).Run()
}
