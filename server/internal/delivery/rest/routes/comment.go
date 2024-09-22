package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
)

func MapCommentRoutes(g *echo.Group, h *handlers.CommentHandlers, mw *middlewares.GuardMiddlewares) {
	g.POST("/:image_id/comments", h.Create(), mw.OnlyAuth)
	g.GET("/:image_id/comments", h.GetByImageID())
	g.PUT("/comments/:comment_id", h.Update(), mw.OnlyAuth)
	g.DELETE("/comments/:comment_id", h.Delete(), mw.OnlyAuth)
	g.GET("/comments/:comment_id/replies", h.GetReplies())

	g.POST("/comments/:comment_id/like", h.LikeComment(), mw.OnlyAuth)
	g.DELETE("/comments/:comment_id/like", h.UnlikeComment(), mw.OnlyAuth)
}
