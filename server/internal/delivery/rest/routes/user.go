package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapUserRoutes(g *echo.Group, h *handlers.UserHandlers, mw *middlewares.AuthMiddlewares) {
	g.GET("/@me", h.Me(), mw.OnlyAuth)
	g.PUT("/:id", h.Update(), mw.OnlyAuth, mw.OwnerOrAdmin)
	g.PUT("/:id/permissions", h.OverwritePermissions(), mw.OnlyAuth, mw.OnlyAdmin)
}
