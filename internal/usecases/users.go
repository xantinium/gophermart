package usecases

import (
	"context"

	"github.com/rs/xid"

	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/tools"
)

func (cases *UseCases) RegisterUser(ctx context.Context, login, password string) error {
	passwordHash, err := tools.HashPassword(password)
	if err != nil {
		return err
	}

	return cases.usersRepo.CreateUser(ctx, login, passwordHash)
}

func (cases *UseCases) AuthorizeUser(ctx context.Context, login, password string) (string, error) {
	user, err := cases.usersRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	hashMatched := tools.CheckPassword(password, user.PasswordHash())
	if !hashMatched {
		return "", models.ErrFailedToMatch
	}

	token := xid.New().String()
	err = cases.tokensRepo.AuthorizeUser(ctx, user.Id(), token)

	return token, err
}

func (cases *UseCases) VerifyUserAuthorization(ctx context.Context, token string) error {
	exists, err := cases.tokensRepo.IsAuthorized(ctx, token)
	if err != nil {
		return err
	}

	if !exists {
		return models.ErrNotFound
	}

	return cases.tokensRepo.RefreshToken(ctx, token)
}

func (cases *UseCases) ClearExpiredTokens(ctx context.Context) error {
	return cases.tokensRepo.ClearExpiredTokens(ctx)
}
