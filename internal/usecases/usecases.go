package usecases

import (
	"context"

	ordersrepo "github.com/xantinium/gophermart/internal/repository/orders"
	tokensrepo "github.com/xantinium/gophermart/internal/repository/tokens"
	usersrepo "github.com/xantinium/gophermart/internal/repository/users"
	withdrawalsrepo "github.com/xantinium/gophermart/internal/repository/withdrawals"
)

type Options struct {
	UsersRepo       *usersrepo.UsersRepository
	TokensRepo      *tokensrepo.TokensRepository
	OrdersRepo      *ordersrepo.OrdersRepository
	WithdrawalsRepo *withdrawalsrepo.WithdrawalsRepository
}

func New(opts Options) *UseCases {
	return &UseCases{
		usersRepo:       opts.UsersRepo,
		tokensRepo:      opts.TokensRepo,
		ordersRepo:      opts.OrdersRepo,
		withdrawalsRepo: opts.WithdrawalsRepo,
	}
}

type UseCases struct {
	usersRepo       *usersrepo.UsersRepository
	tokensRepo      *tokensrepo.TokensRepository
	ordersRepo      *ordersrepo.OrdersRepository
	withdrawalsRepo *withdrawalsrepo.WithdrawalsRepository
}

type Balance struct {
	AvaliableAccrual float32
	TotalWithdrawn   float32
}

func (cases *UseCases) GetUserBalance(ctx context.Context, userID int) (Balance, error) {
	var (
		err                          error
		totalAccrual, totalWithdrawn float32
	)

	totalAccrual, err = cases.ordersRepo.GetTotalAccrual(ctx, userID)
	if err != nil {
		return Balance{}, err
	}

	totalWithdrawn, err = cases.withdrawalsRepo.GetTotalWithdrawn(ctx, userID)
	if err != nil {
		return Balance{}, err
	}

	return Balance{
		AvaliableAccrual: totalAccrual - totalWithdrawn,
		TotalWithdrawn:   totalWithdrawn,
	}, err
}
