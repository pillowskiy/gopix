package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthGoogleProvider struct {
	cfg *oauth2.Config
}

func NewOAuthGoogleProvider(cfg *config.OAuthGoogle) *OAuthGoogleProvider {
	oAuthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &OAuthGoogleProvider{cfg: oAuthCfg}
}

func (p *OAuthGoogleProvider) GetUserInfo(ctx context.Context, code string) (*domain.OAuthUser, error) {
	token, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, ErrIncorrectCode
	}

	client := p.cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	oauthUser := &domain.OAuthUser{Service: domain.OAuthServiceGoogle}
	if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return oauthUser, nil
}

func (p *OAuthGoogleProvider) GetAuthURL() string {
	// WARNING: Unsafe. Add CSRF Protection
	return p.cfg.AuthCodeURL("")
}
