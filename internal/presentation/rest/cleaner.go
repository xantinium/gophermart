package rest

import (
	"context"
	"sync"
	"time"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/usecases"
)

const expiredTokensCheckInterval = time.Second * 10

func NewTokensCleaner(cases *usecases.UseCases) *TokensCleaner {
	return &TokensCleaner{
		cases: cases,
	}
}

type TokensCleaner struct {
	wg    sync.WaitGroup
	cases *usecases.UseCases
}

func (cleaner *TokensCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(expiredTokensCheckInterval)

	cleaner.wg.Add(1)
	go func() {
		defer cleaner.wg.Done()

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				err := cleaner.cases.ClearExpiredTokens(ctx)
				if err != nil {
					logger.Errorf("failed to clear expired tokens: %v", err)
				}
			}
		}
	}()
}

func (cleaner *TokensCleaner) Wait() {
	cleaner.wg.Wait()
}
