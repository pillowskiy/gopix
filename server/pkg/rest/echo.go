package rest

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type ReqIDCtxKey struct{}
type UserCtxKey struct{}

func GetEchoRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetEchoCtxWithRequestID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetEchoRequestID(c))
	return ctx, cancel
}

func GetEchoRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetEchoRequestID(c))
}

func GetEchoIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

func DecodeEchoBody(c echo.Context, s interface{}) error {
	if err := c.Bind(s); err != nil {
		return err
	}
	return nil
}
