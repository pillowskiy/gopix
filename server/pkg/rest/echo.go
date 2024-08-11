package rest

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
)

type echoContextAdapter struct {
	c echo.Context
}

func NewEchoContextAdapter(c echo.Context) *echoContextAdapter {
	return &echoContextAdapter{c: c}
}

func MakeEchoHandler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(NewEchoContextAdapter(c))
	}
}

func (a *echoContextAdapter) GetRequestID() string {
	return a.c.Response().Header().Get(echo.HeaderXRequestID)
}

func (a *echoContextAdapter) GetCtxWithRequestID() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(a.c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, a.GetRequestID())
	return ctx, cancel
}
func (a *echoContextAdapter) GetRequestCtx() context.Context {
	return context.WithValue(a.c.Request().Context(), ReqIDCtxKey{}, a.GetRequestID())
}

func (a *echoContextAdapter) GetIPAddress() string {
	return a.c.Request().RemoteAddr
}

func (a *echoContextAdapter) DecodeBody(s interface{}) error {
	if err := a.c.Bind(s); err != nil {
		return err
	}
	return nil
}

func (a *echoContextAdapter) Cookie(name string) (*http.Cookie, error) {
	return a.c.Cookie(name)
}

func (a *echoContextAdapter) SetCookie(cookie *http.Cookie) {
	a.c.SetCookie(cookie)
}

func (a *echoContextAdapter) Cookies() []*http.Cookie {
	return a.c.Cookies()
}

func (a *echoContextAdapter) WriteJSON(status int, v interface{}) error {
	return a.c.JSON(status, v)
}

func (a *echoContextAdapter) NoContent(status int) error {
	return a.c.NoContent(status)
}

func (a *echoContextAdapter) Redirect(status int, url string) error {
	return a.c.Redirect(status, url)
}

func (a *echoContextAdapter) Path() string {
	return a.c.Path()
}

func (a *echoContextAdapter) SetPath(p string) {
	a.c.SetPath(p)
}

func (a *echoContextAdapter) Param(name string) string {
	return a.c.Param(name)
}

func (a *echoContextAdapter) ParamNames() []string {
	return a.c.ParamNames()
}

func (a *echoContextAdapter) SetParamNames(names ...string) {
	a.c.SetParamNames(names...)
}

func (a *echoContextAdapter) ParamValues() []string {
	return a.c.ParamValues()
}

func (a *echoContextAdapter) SetParamValues(values ...string) {
	a.c.SetParamValues(values...)
}

func (a *echoContextAdapter) QueryParam(name string) string {
	return a.c.QueryParam(name)
}

func (a *echoContextAdapter) QueryParams() url.Values {
	return a.c.QueryParams()
}

func (a *echoContextAdapter) QueryString() string {
	return a.c.QueryString()
}
