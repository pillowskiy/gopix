package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/storage"
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

	s.echo.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Pong!"})
	})

	s.logger.Infof("Server is listening on ADDR: %s", s.cfg.Addr)
	if err := s.echo.StartServer(server); err != nil {
		return err
	}

	return nil
}
