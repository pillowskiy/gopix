package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapAlbumRoutes(g *echo.Group, h *handlers.AlbumHandlers, mw *middlewares.GuardMiddlewares) {
	g.POST("/", h.Create(), mw.OnlyAuth)
	g.GET("/users/:user_id", h.GetByAuthorID())
	g.DELETE("/:album_id", h.Delete(), mw.OnlyAuth)
	g.PUT("/:album_id", h.Update(), mw.OnlyAuth)

	g.POST("/:album_id/images/:image_id", h.PutImage(), mw.OnlyAuth)
	g.DELETE("/:album_id/images/:image_id", h.DeleteImage(), mw.OnlyAuth)
	g.GET("/:album_id/images", h.GetAlbumImages())
}
