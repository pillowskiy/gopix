package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type authUseCase interface {
	Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
}

type AuthHandlers struct {
	uc     authUseCase
	logger logger.Logger
}

func NewAuthHandlers(uc authUseCase, logger logger.Logger) *AuthHandlers {
	return &AuthHandlers{uc: uc, logger: logger}
}

type registerDTO struct {
	Username string `json:"username" validate:"required,gte=6,lte=60"`
	Email    string `json:"email" validate:"required,lte=60,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (h *AuthHandlers) Register(c rest.Context) error {
	ctx := c.GetRequestCtx()

	reg := new(registerDTO)
	if err := c.DecodeBody(reg); err != nil {
		h.logger.Errorf("Register.DecodeBody: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := validator.ValidateStruct(ctx, reg); err != nil {
		return c.WriteJSON(rest.NewBadRequestError("Registration body has incorrect type").Response())
	}

	user := &domain.User{
		Username:     reg.Username,
		PasswordHash: reg.Password,
		Email:        reg.Email,
	}

	newUser, err := h.uc.Register(ctx, user)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyExists) {
			return c.WriteJSON(rest.NewConflictError("User already exists").Response())
		}
		return c.WriteJSON(rest.NewInternalServerError().Response())
	}

	return c.WriteJSON(http.StatusCreated, newUser)
}
