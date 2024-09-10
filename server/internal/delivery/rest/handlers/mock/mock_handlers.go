package mock_handlers

import (
	context "context"

	"github.com/labstack/echo/v4"
	domain "github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/rest"
)

func NewMockCtxUser() (*domain.User, func(c echo.Context)) {
	ctxUser := &domain.User{
		ID:          1,
		Username:    "username",
		Email:       "username@gmail.com",
		Permissions: 1,
		AvatarURL:   "https://example.com/username.png",
	}

	mockCtxUser := func(c echo.Context) {
		c.Set("user", ctxUser)
		ctx := context.WithValue(c.Request().Context(), rest.UserCtxKey{}, ctxUser)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	return ctxUser, mockCtxUser
}
