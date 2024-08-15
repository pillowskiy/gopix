package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/pkg/rest"
)

func MapEchoAuthRoutes(g *echo.Group, h *handlers.AuthHandlers) {
	g.POST("/register", rest.MakeEchoHandler(h.Register))
	g.POST("/login", rest.MakeEchoHandler(h.Login))
}
