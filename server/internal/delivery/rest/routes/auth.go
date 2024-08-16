package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapAuthRoutes(g *echo.Group, h *handlers.AuthHandlers, mw *middlewares.AuthMiddlewares) {
	g.POST("/register", h.Register())
	g.POST("/login", h.Login())
	g.DELETE("/logout", h.Logout(), mw.OnlyAuth)
}
