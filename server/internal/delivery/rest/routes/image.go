package routes

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
	"github.com/pillowskiy/gopix/internal/domain"
)

func MapImageRoutes(g *echo.Group, h *handlers.ImageHandlers, mw *middlewares.GuardMiddlewares) {
	g.POST("/",
		h.Create(),
		mw.OnlyAuth,
		mw.WithSomePermission(domain.PermissionsUploadImage),
		middlewares.TimeoutMiddleware(15*time.Minute),
	)
}
