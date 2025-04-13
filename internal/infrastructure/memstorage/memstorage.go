package memstorage

import (
	"context"
	"sync"
	"time"

	"github.com/xantinium/gophermart/internal/models"
)

const tokenLifeTime = 24 * time.Hour

func New() *MemStorage {
	return &MemStorage{
		tokens: make(map[string]tokenInfo),
	}
}

type MemStorage struct {
	mx     sync.RWMutex
	tokens map[string]tokenInfo
}

func (storage *MemStorage) GetByToken(_ context.Context, token string) (int, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	info, ok := storage.tokens[token]
	if !ok {
		return 0, models.ErrNotFound
	}

	return info.UserID, nil
}

func (storage *MemStorage) SetToken(_ context.Context, userID int, token string) error {
	now := time.Now()

	storage.mx.Lock()
	defer storage.mx.Unlock()

	for token := range storage.tokens {
		if storage.tokens[token].UserID == userID {
			delete(storage.tokens, token)
			break
		}
	}

	storage.tokens[token] = tokenInfo{
		UserID:         userID,
		ExpirationTime: now.Add(tokenLifeTime),
	}

	return nil
}

func (storage *MemStorage) RefreshToken(_ context.Context, token string) error {
	now := time.Now()

	storage.mx.Lock()
	defer storage.mx.Unlock()

	info, ok := storage.tokens[token]
	if !ok {
		return models.ErrNotFound
	}

	storage.tokens[token] = tokenInfo{
		UserID:         info.UserID,
		ExpirationTime: now.Add(tokenLifeTime),
	}

	return nil
}

func (storage *MemStorage) ClearExpiredTokens(_ context.Context) error {
	now := time.Now()

	storage.mx.Lock()
	defer storage.mx.Unlock()

	nonExpiredTokens := make(map[string]tokenInfo)
	for token := range storage.tokens {
		if now.Before(storage.tokens[token].ExpirationTime) {
			nonExpiredTokens[token] = storage.tokens[token]
		}
	}

	storage.tokens = nonExpiredTokens

	return nil
}

type tokenInfo struct {
	UserID         int
	ExpirationTime time.Time
}
