package oauth

import (
	"context"
	"errors"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
)

var (
	ErrIncorrectCode = errors.New("incorrect code")
	UnknownService   = errors.New("unknown service")
)

type OAuthProvider interface {
	GetUserInfo(ctx context.Context, code string) (*domain.OAuthUser, error)
	GetAuthURL() string
}

type OAuthClient struct {
	cfg    *config.OAuth
	google OAuthProvider
}

func NewOAuthClient(cfg *config.OAuth) *OAuthClient {
	return &OAuthClient{
		cfg:    cfg,
		google: NewOAuthGoogleProvider(cfg.Google),
	}
}

func (client *OAuthClient) GetUserInfo(ctx context.Context, code string, service domain.OAuthService) (*domain.OAuthUser, error) {
	switch service {
	case domain.OAuthServiceGoogle:
		return client.google.GetUserInfo(ctx, code)
	default:
		return nil, UnknownService
	}
}

func (client *OAuthClient) GetAuthURL(service domain.OAuthService) (string, error) {
	switch service {
	case domain.OAuthServiceGoogle:
		return client.google.GetAuthURL(), nil
	default:
		return "", UnknownService
	}
}
