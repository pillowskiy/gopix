package routes

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
	"github.com/pillowskiy/gopix/internal/domain"
)

func MapImageRoutes(g *echo.Group, h *handlers.ImageHandlers, mw *middlewares.GuardMiddlewares) {
	g.GET("/", h.GetDiscover())
	g.GET("/favorites/:user_id", h.Favorites())

	g.POST("/",
		h.Upload(),
		mw.OnlyAuth,
		mw.WithSomePermission(domain.PermissionsUploadImage),
		middlewares.TimeoutMiddleware(15*time.Minute),
	)

	g.DELETE("/:id", h.Delete(), mw.OnlyAuth)
	g.PUT("/:id", h.Update(), mw.OnlyAuth)
	g.GET("/:id", h.GetDetailed(), mw.OptionalAuth)
	g.GET("/:id/similar", h.Similar())

	g.GET("/:id/states", h.GetStates(), mw.OnlyAuth)

	g.POST("/:id/like", h.AddLike(), mw.OnlyAuth)
	g.DELETE("/:id/like", h.RemoveLike(), mw.OnlyAuth)
}
