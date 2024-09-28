package httprepo

import (
	"context"
	"fmt"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
)

type OAuthProvider interface {
	GetUserInfo(ctx context.Context, code string) (*domain.OAuthUser, error)
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
		return nil, fmt.Errorf("OAuthRepository.GetUserInfo: invalid oauth service %s", service)
	}
}
