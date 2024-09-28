package httprepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/domain"
	repository "github.com/pillowskiy/gopix/internal/respository"
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
	fmt.Printf("Google OAuth URL: %s", oAuthCfg.AuthCodeURL(""))
	return &OAuthGoogleProvider{cfg: oAuthCfg}
}

func (o *OAuthGoogleProvider) GetUserInfo(ctx context.Context, code string) (*domain.OAuthUser, error) {
	token, err := o.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, repository.ErrIncorrectInput
	}

	client := o.cfg.Client(ctx, token)
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
