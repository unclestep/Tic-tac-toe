package di

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
	httpDelivery "tictactoe/internal/delivery/http"
	"tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/handler"
	"tictactoe/internal/delivery/http/json/middleware"
	"tictactoe/internal/domain/service"

	// "tictactoe/internal/infrastructure/storage/memory"
	"tictactoe/internal/infrastructure/storage/ds"
	"tictactoe/internal/infrastructure/storage/postgres"
	"tictactoe/internal/infrastructure/storage/repository"

	"go.uber.org/fx"
)

var TicTacToe = fx.Module(
	"Tic-Tac-Toe",
	domain,
	app,
	repo,
	dataSource,
	delivery,
	fx.Invoke(registerServer),
)

func registerServer(lc fx.Lifecycle, router *httpDelivery.Router) {
	server := http.Server{
		Handler: middleware.WithRecovery(router.Handler()),
		Addr:    os.Getenv("TICTACTOE_ADDR"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server error: %s", err)
				}
			}()
			log.Printf("\nserver has started on: %s", os.Getenv("TICTACTOE_ADDR"))
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
		fx.As(new(port.BotMover)),
	)),
	fx.Provide(fx.Annotate(
		service.NewHumanMovement,
		fx.As(new(port.HumanMover)),
	)),
	fx.Provide(fx.Annotate(
		service.NewUserService,
		fx.As(new(port.UserService)),
	)),
	fx.Provide(service.NewWinChecker),
	fx.Provide(service.NewHeuristic),
)

var app = fx.Module(
	"Application",
	fx.Provide(fx.Annotate(
		usecase.NewConnect,
		fx.As(new(port.ConnectUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewCreate,
		fx.As(new(port.CreateUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewStart,
		fx.As(new(port.StartUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewMakeMove,
		fx.As(new(port.MakeMoveUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewDisconnect,
		fx.As(new(port.DisconnectUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewSignUp,
		fx.As(new(usecase.SignUpUseCase)),
	)),
	fx.Provide(fx.Annotate(
		usecase.NewSignIn,
		fx.As(new(usecase.SignInUseCase)),
	)),
)

var repo = fx.Module(
	"Repository",
	fx.Provide(fx.Annotate(
		repository.NewRulesRepo,
		fx.As(new(port.RulesRepo)),
	)),
	fx.Provide(fx.Annotate(
		repository.NewSessionRepo,
		fx.As(new(port.SessionRepo)),
	)),
	fx.Provide(fx.Annotate(
		repository.NewUserRepo,
		fx.As(new(port.UserRepo)),
	)),
)

var dataSource = fx.Module(
	"DataSource",
	fx.Provide(func() (postgres.DBTX, error) {
		return postgres.NewPool(context.Background(), os.Getenv("POSTGRES_DSN"))
	}),
	fx.Provide(fx.Annotate(
		postgres.NewRulesDataSource,
		fx.As(new(ds.RulesDataSource)),
	)),
	fx.Provide(fx.Annotate(
		postgres.NewSessionDataSource,
		fx.As(new(ds.SessionDataSource)),
	)),
	fx.Provide(fx.Annotate(
		postgres.NewUserDataSource,
		fx.As(new(ds.UserDataSource)),
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
		handler.NewSignUpHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"signup_h"`),
	)),
	fx.Provide(fx.Annotate(
		handler.NewSignInHandler,
		fx.As(new(http.Handler)),
		fx.ResultTags(`name:"signin_h"`),
	)),
	fx.Provide(middleware.NewUserAuthenticator),
	fx.Provide(fx.Annotate(
		httpDelivery.NewRouter,
		fx.ParamTags(
			`name:"create_h"`,
			`name:"connect_h"`,
			`name:"start_h"`,
			`name:"makemove_h"`,
			`name:"disconnect_h"`,
			`name:"signup_h"`,
			`name:"signin_h"`,
		),
	)),
	fx.Provide(json.NewErrorMapper),
)
