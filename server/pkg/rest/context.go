package rest

import (
	"context"
	"net/http"
	"net/url"
)

// Is this overhead - absolutely, but why not
type Context interface {
	GetRequestID() string
	GetCtxWithRequestID() (context.Context, context.CancelFunc)
	GetRequestCtx() context.Context
	GetIPAddress() string

	DecodeBody(s interface{}) error

	Cookie(name string) (*http.Cookie, error)
	SetCookie(cookie *http.Cookie)
	Cookies() []*http.Cookie

	WriteJSON(status int, v interface{}) error
	NoContent(status int) error
	Redirect(code int, url string) error

	Path() string
	SetPath(p string)

	Param(name string) string
	ParamNames() []string
	SetParamNames(names ...string)
	ParamValues() []string
	SetParamValues(values ...string)

	QueryParam(name string) string
	QueryParams() url.Values
	QueryString() string
}

type ReqIDCtxKey struct{}

type HandlerFunc func(c Context) error
