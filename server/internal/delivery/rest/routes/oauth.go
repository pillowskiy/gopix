package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/domain"
)

func MapOAuthRoutes(g *echo.Group, h *handlers.OAuthHandlers) {
	g.GET("/callback/google", h.Callback(domain.OAuthServiceGoogle))
	g.GET("/google", h.AuthRedirect(domain.OAuthServiceGoogle))
}
