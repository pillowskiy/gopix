package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapNotificationRoutes(g *echo.Group, h *handlers.NotificationHandlers, mw *middlewares.GuardMiddlewares) {
	g.GET("/", h.GetNotifications()) // mw.OnlyAuth)
	g.GET("/sse", h.NotificationsSSE())
	g.GET("/stats", h.GetStats())
	g.POST("/", h.Notify(), mw.OnlyAuth, mw.OnlyAdmin)
}
