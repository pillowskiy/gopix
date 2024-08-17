package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
	"github.com/pillowskiy/gopix/internal/delivery/rest/routes"
	"github.com/pillowskiy/gopix/internal/respository/postgres"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/storage"
	"github.com/pillowskiy/gopix/pkg/token"
)

type EchoServer struct {
	echo   *echo.Echo
	cfg    *config.Server
	sh     *storage.StorageHolder
	logger logger.Logger
}

func NewEchoServer(cfg *config.Server, sh *storage.StorageHolder, logger logger.Logger) *EchoServer {
	return &EchoServer{echo: echo.New(), cfg: cfg, sh: sh, logger: logger}
}

func (s *EchoServer) Listen() error {
	server := &http.Server{
		Addr:         s.cfg.Addr,
		ReadTimeout:  time.Second * s.cfg.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.WriteTimeout,
	}

	if err := s.MapHandlers(); err != nil {
		return err
	}

	s.logger.Infof("Server is listening on ADDR: %s", s.cfg.Addr)
	if err := s.echo.StartServer(server); err != nil {
		return err
	}

	return nil
}

func (s *EchoServer) MapHandlers() error {
	userRepo := postgres.NewUserRepository(s.sh.Postgres)

	jwtTokenGen := token.NewJWTTokenGenerator(
		s.cfg.Session.Secret,
		s.cfg.Session.Expire*time.Second,
	)
	authUC := usecase.NewAuthUseCase(userRepo, s.logger, jwtTokenGen)
	userUC := usecase.NewUserUseCase(userRepo)

	v1 := s.echo.Group("/api/v1")
	authMiddlewares := middlewares.NewAuthMiddlewares(authUC, s.logger, s.cfg.Cookie)

	authGroup := v1.Group("/auth")
	authHandlers := handlers.NewAuthHandlers(authUC, s.logger, s.cfg.Cookie)
	routes.MapAuthRoutes(authGroup, authHandlers, authMiddlewares)

	userGroup := v1.Group("/users")
	userHandlers := handlers.NewUserHandlers(userUC, s.logger)
	routes.MapUserRoutes(userGroup, userHandlers, authMiddlewares)

	return nil
}
