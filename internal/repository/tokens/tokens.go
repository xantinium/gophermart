package tokensrepo

import "context"

func New(storage TokensStorage) *TokensRepository {
	return &TokensRepository{
		storage: storage,
	}
}

type TokensRepository struct {
	storage TokensStorage
}

func (repo *TokensRepository) GetAuthorizedUser(ctx context.Context, token string) (int, error) {
	return repo.storage.GetByToken(ctx, token)
}

func (repo *TokensRepository) AuthorizeUser(ctx context.Context, userID int, token string) error {
	return repo.storage.SetToken(ctx, userID, token)
}

func (repo *TokensRepository) RefreshToken(ctx context.Context, token string) error {
	return repo.storage.RefreshToken(ctx, token)
}

func (repo *TokensRepository) ClearExpiredTokens(ctx context.Context) error {
	return repo.storage.ClearExpiredTokens(ctx)
}
