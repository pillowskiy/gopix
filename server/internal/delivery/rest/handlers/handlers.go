package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
)

func GetContextUser(c echo.Context) (*domain.User, error) {
	user, ok := c.Get("user").(*domain.User)
	if !ok || user == nil {
		return nil, errors.New("cannot get user from context")
	}
	return user, nil
}
