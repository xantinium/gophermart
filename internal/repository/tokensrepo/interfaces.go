package tokensrepo

import "context"

type TokensStorage interface {
	HasToken(ctx context.Context, token string) (bool, error)
	SetToken(ctx context.Context, userID int, token string) error
	RefreshToken(ctx context.Context, token string) error
	ClearExpiredTokens(ctx context.Context) error
}
