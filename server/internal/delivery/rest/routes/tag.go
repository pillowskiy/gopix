package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapTagRoutes(g *echo.Group, h *handlers.TagHandlers, mw *middlewares.GuardMiddlewares) {
	g.GET("/tags", h.Search(), mw.OnlyAuth)
	g.POST("/tags", h.Create(), mw.OnlyAuth, mw.OnlyAdmin)
	g.DELETE("/tags/:tag_id", h.Delete(), mw.OnlyAuth, mw.OnlyAdmin)

	g.PUT("/:image_id/tags", h.UpsertImageTag(), mw.OnlyAuth)
	g.DELETE("/:image_id/tags/:tag_id", h.DeleteImageTag(), mw.OnlyAuth)
}
