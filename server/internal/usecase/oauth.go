package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/infrastructure/oauth"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pkg/errors"
)

type OAuthRepository interface {
	GetByOAuthID(ctx context.Context, oauthID string) (*domain.OAuth, error)
	Create(ctx context.Context, oauth *domain.OAuth) error

	repository.Transactional
}

type OAuthAuthUC interface {
	GetUnique(ctx context.Context, user *domain.User) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
	GenerateToken(user *domain.User) (string, error)
}

type OAuthClient interface {
	GetUserInfo(ctx context.Context, code string, sevice domain.OAuthService) (*domain.OAuthUser, error)
	GetAuthURL(service domain.OAuthService) (string, error)
}

type oAuthUseCase struct {
	repo   OAuthRepository
	authUC OAuthAuthUC
	client OAuthClient
}

func NewOAuthUseCase(repo OAuthRepository, authUC OAuthAuthUC, client OAuthClient) *oAuthUseCase {
	return &oAuthUseCase{repo: repo, authUC: authUC, client: client}
}

func (uc *oAuthUseCase) GetAuthURL(ctx context.Context, service domain.OAuthService) (string, error) {
	return uc.client.GetAuthURL(service)
}

func (uc *oAuthUseCase) Authenticate(ctx context.Context, code string, service domain.OAuthService) (*domain.UserWithToken, error) {
	extUser, err := uc.client.GetUserInfo(ctx, code, service)
	if err != nil {
		if errors.Is(err, oauth.ErrIncorrectCode) {
			return nil, ErrUnprocessable
		}
		return nil, errors.Wrap(err, "OAuthUseCase.Authenticate.GetUserInfo")
	}

	oauth, err := uc.GetByOAuthID(ctx, extUser.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return uc.register(ctx, extUser)
		}
		return nil, errors.Wrap(err, "OAuthUseCase.Authenticate.GetByOAuthID")
	}

	if oauth.Service != extUser.Service {
		return nil, fmt.Errorf("OAuthUseCase.Authenticate: invalid oauth service %s", oauth.Service)
	}

	return uc.login(ctx, extUser, oauth)
}

func (uc *oAuthUseCase) register(ctx context.Context, oauthUser *domain.OAuthUser) (*domain.UserWithToken, error) {
	uniqueUsername := fmt.Sprintf("%s%s", strings.Split(oauthUser.Email, "@")[0], oauthUser.ID[len(oauthUser.ID)-5:])
	userInput := &domain.User{
		Username:  uniqueUsername,
		Email:     oauthUser.Email,
		AvatarURL: oauthUser.Picture,
		External:  true,
	}

	authUser := new(domain.UserWithToken)
	err := uc.repo.DoInTransaction(ctx, func(ctx context.Context) error {
		createdUser, err := uc.authUC.Register(ctx, userInput)
		if err != nil {
			return err
		}

		oauthInput := &domain.OAuth{
			OAuthID: oauthUser.ID,
			UserID:  createdUser.User.ID,
			Service: oauthUser.Service,
		}
		if err := uc.repo.Create(ctx, oauthInput); err != nil {
			return err
		}

		authUser = createdUser
		return nil
	})

	return authUser, err
}

func (uc *oAuthUseCase) login(ctx context.Context, oauthUser *domain.OAuthUser, oauth *domain.OAuth) (*domain.UserWithToken, error) {
	user, err := uc.authUC.GetUnique(ctx, &domain.User{
		Email: oauthUser.Email,
	})
	if err != nil {
		return nil, err
	}

	if oauth.UserID != user.ID {
		return nil, ErrInvalidCredentials
	}

	if !user.External {
		return nil, fmt.Errorf("oauth referred to non-external user (id: %v)", user.ID)
	}

	t, err := uc.authUC.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithToken{
		User:  user,
		Token: t,
	}, nil
}

func (uc *oAuthUseCase) GetByOAuthID(ctx context.Context, oauthID string) (*domain.OAuth, error) {
	oauth, err := uc.repo.GetByOAuthID(ctx, oauthID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "OAuthUseCase.GetByOAuthID")
	}

	return oauth, nil
}
