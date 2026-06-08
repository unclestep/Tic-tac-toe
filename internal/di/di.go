package di

import (
	"context"
	"log"
	"net/http"
	"os"

	aport "tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
	httpDelivery "tictactoe/internal/delivery/http"
	"tictactoe/internal/delivery/http/json/handler"
	"tictactoe/internal/domain/service"
	"tictactoe/internal/infrastructure/storage/memory"
	sport "tictactoe/internal/infrastructure/storage/port"
	"tictactoe/internal/infrastructure/storage/repository"

	"go.uber.org/fx"
)

var TicTacToe = fx.Module(
	"Tic-Tac-Toe",
	domain,
	app,
	repo,
	ds,
	delivery,
	fx.Invoke(registerServer),
)

func registerServer(lc fx.Lifecycle, router *httpDelivery.Router) {
	server := http.Server{
		Handler: router.Handler(),
		Addr:    "localhost:12121",
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil {
					log.Fatalf("server error: %s", err)
				}
			}()
			log.Printf("\nserver has started on: %s", os.Getenv("SERVER_ADDR"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

var domain = fx.Module(
	"Domain",
	fx.Provide(fx.Annotate(
		service.NewBotMovement,
		fx.As(new(usecase.BotMover)),
	)),
	fx.Provide(fx.Annotate(
		service.NewHumanMovement,
		fx.As(new(usecase.HumanMover)),
	)),
	fx.Provide(service.NewWinChecker),
	fx.Provide(service.NewHeuristic),
)

var app = fx.Module(
	"Application",
	fx.Provide(fx.Annotate(
		usecase.NewConnect,
		fx.As(new(aport.ConnectUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewCreate,
		fx.As(new(aport.CreateUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewStart,
		fx.As(new(aport.StartUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewMakeMove,
		fx.As(new(aport.MakeMoveUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewDisconnect,
		fx.As(new(aport.DisconnectUseCase)),
	)),
)

var repo = fx.Module(
	"Repository",
	fx.Provide(fx.Annotate(
		repository.NewRulesRepo,
		fx.As(new(aport.RulesRepo)),
	)),
	fx.Provide(fx.Annotate(
		repository.NewSessionRepo,
		fx.As(new(aport.SessionRepo)),
	)),
)

var ds = fx.Module(
	"DataSource",
	fx.Provide(fx.Annotate(
		memory.NewRulesDataSource,
		fx.As(new(sport.RulesDataSource)),
	)),
	fx.Provide(fx.Annotate(
		memory.NewSessionDataSource,
		fx.As(new(sport.SessionDataSource)),
	)),
)

var delivery = fx.Module(
	"HTTP",
	fx.Provide(fx.Annotate(
		handler.NewConnectHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"connect_h"`),
	)),
	fx.Provide(fx.Annotate(
		handler.NewCreateHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"create_h"`),
	)),
	fx.Provide(fx.Annotate(
		handler.NewStartHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"start_h"`),
	)),
	fx.Provide(fx.Annotate(
		handler.NewMakeMoveHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"makemove_h"`),
	)),
	fx.Provide(fx.Annotate(
		handler.NewDisconnectHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"disconnect_h"`),
	)),
	fx.Provide(fx.Annotate(
		httpDelivery.NewRouter,
		fx.ParamTags(
			`name:"create_h"`,
			`name:"connect_h"`,
			`name:"start_h"`,
			`name:"makemove_h"`,
			`name:"disconnect_h"`,
		),
	)),
	fx.Provide(handler.NewErrorMapper),
)
