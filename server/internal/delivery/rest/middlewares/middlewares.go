package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func TimeoutMiddleware(
	timeout time.Duration,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			c.SetRequest(c.Request().WithContext(ctx))
			next(c)

			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.Println("HTTP Request timed out")
					return c.JSON(http.StatusRequestTimeout, "Request timed out")
				}
			}

			return nil
		}
	}
}
