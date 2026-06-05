package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"

	"tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
	httpRouter "tictactoe/internal/delivery/http"
	"tictactoe/internal/delivery/http/json/handler"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/service"
	servicePort "tictactoe/internal/domain/service/port"
	"tictactoe/internal/infrastructure/storage/memory"
	storagePort "tictactoe/internal/infrastructure/storage/port"
	"tictactoe/internal/infrastructure/storage/repository"
)

var Module = fx.Options(
	infraModule,
	domainModule,
	appModule,
	deliveryModule,
)

var infraModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			memory.NewMemorySessionDataSource,
			fx.As(new(storagePort.SessionDataSource)),
		),
	),
	fx.Provide(
		fx.Annotate(
			memory.NewMemoryRulesDataSource,
			fx.As(new(storagePort.RulesDataSource)),
		),
	),

	fx.Provide(
		fx.Annotate(
			repository.NewSessionRepo,
			fx.As(new(port.SessionRepo)),
		),
	),
	fx.Provide(
		fx.Annotate(
			repository.NewRulesRepo,
			fx.As(new(port.RulesRepo)),
		),
	),
)

var domainModule = fx.Options(
	fx.Provide(model.NewDefaultRules),

	fx.Provide(service.NewWinCheckerService),
	fx.Provide(service.NewAdvancer),
	fx.Provide(service.NewBotLogic),
	fx.Provide(service.NewMovementService),
	fx.Provide(service.NewGamePreparer),

	fx.Provide(
		fx.Annotate(
			service.NewGameMechanics,
			fx.As(new(servicePort.GameMechanics)),
		),
	),
)

var appModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			usecase.NewCreate,
			fx.As(new(port.CreateUseCase)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewConnect,
			fx.As(new(port.ConnectUseCase)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewStart,
			fx.As(new(port.StartUseCase)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewMakeMove,
			fx.As(new(port.MakeMoveUseCase)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewDisconnect,
			fx.As(new(port.DisconnectUseCase)),
		),
	),
)

var deliveryModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			handler.NewCreateHandler,
			fx.As(new(http.Handler)),
			fx.ResultTags(`name:"h_create"`),
		),
		fx.Annotate(
			handler.NewConnectHandler,
			fx.As(new(http.Handler)),
			fx.ResultTags(`name:"h_connect"`),
		),
		fx.Annotate(
			handler.NewStartHandler,
			fx.As(new(http.Handler)),
			fx.ResultTags(`name:"h_start"`),
		),
		fx.Annotate(
			handler.NewMakeMoveHandler,
			fx.As(new(http.Handler)),
			fx.ResultTags(`name:"h_move"`),
		),
		fx.Annotate(
			handler.NewDisconnectHandler,
			fx.As(new(http.Handler)),
			fx.ResultTags(`name:"h_disconnect"`),
		),
	),

	fx.Provide(
		fx.Annotate(
			httpRouter.NewRouter,
			fx.ParamTags(
				`name:"h_create"`,
				`name:"h_connect"`,
				`name:"h_start"`,
				`name:"h_move"`,
				`name:"h_disconnect"`,
			),
		),
	),

	fx.Invoke(startServer),
)

func startServer(lc fx.Lifecycle, router *httpRouter.Router) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal(err)
				}
			}()
			fmt.Println("Server started on :8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
