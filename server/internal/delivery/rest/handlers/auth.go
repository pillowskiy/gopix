package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/rest"
	"github.com/pillowskiy/gopix/pkg/validator"
)

type authUseCase interface {
	Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
	Login(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
}

type AuthHandlers struct {
	uc     authUseCase
	logger logger.Logger
	cfg    *config.Cookie
}

func NewAuthHandlers(uc authUseCase, logger logger.Logger, cfg *config.Cookie) *AuthHandlers {
	return &AuthHandlers{uc: uc, logger: logger, cfg: cfg}
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
		return c.WriteJSON(rest.NewBadRequestError("Registration body has incorrect type").Response())
	}

	if err := validator.ValidateStruct(ctx, reg); err != nil {
		return c.WriteJSON(rest.NewBadRequestError("Registration body has incorrect type").Response())
	}

	user := &domain.User{
		Username:     reg.Username,
		PasswordHash: reg.Password,
		Email:        reg.Email,
	}

	authUser, err := h.uc.Register(ctx, user)
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyExists) {
			return c.WriteJSON(rest.NewConflictError("User already exists").Response())
		}
		return c.WriteJSON(rest.NewInternalServerError().Response())
	}

	h.storeToken(c, authUser.Token)
	return c.WriteJSON(http.StatusCreated, authUser)
}

type loginDTO struct {
	Initials string `json:"initials" validate:"required,gte=2"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (h *AuthHandlers) Login(c rest.Context) error {
	ctx := c.GetRequestCtx()

	login := new(loginDTO)
	if err := c.DecodeBody(login); err != nil {
		h.logger.Errorf("Login.DecodeBody: %v", err)
		return c.WriteJSON(rest.NewBadRequestError("Login body has incorrect type").Response())
	}

	if err := validator.ValidateStruct(ctx, login); err != nil {
		return c.WriteJSON(rest.NewBadRequestError("Login body has incorrect type").Response())
	}

	credentials := &domain.User{
		Username:     login.Initials,
		PasswordHash: login.Password,
	}

	authUser, err := h.uc.Login(ctx, credentials)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return c.WriteJSON(rest.NewUnauthorizedError("Invalid credentials").Response())
		}
		return c.WriteJSON(rest.NewInternalServerError().Response())
	}

	h.storeToken(c, authUser.Token)
	return c.WriteJSON(http.StatusOK, authUser)
}

func (h *AuthHandlers) Logout(c rest.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	return c.WriteJSON(http.StatusOK, true)
}

func (h *AuthHandlers) storeToken(c rest.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     h.cfg.Name,
		Value:    token,
		Path:     "/",
		HttpOnly: h.cfg.HttpOnly,
		MaxAge:   int((h.cfg.Expire * time.Second).Seconds()),

		SameSite: http.SameSiteLaxMode,
	})
}
