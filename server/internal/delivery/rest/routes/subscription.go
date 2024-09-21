package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapSubscriptionRoutes(g *echo.Group, h *handlers.SubscriptionHandlers, mw *middlewares.GuardMiddlewares) {
	g.POST("/:user_id", h.Follow(), mw.OnlyAuth)
	g.DELETE("/:user_id", h.Unfollow(), mw.OnlyAuth)
}
