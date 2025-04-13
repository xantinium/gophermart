package usecases

import (
	ordersrepo "github.com/xantinium/gophermart/internal/repository/orders"
	tokensrepo "github.com/xantinium/gophermart/internal/repository/tokens"
	usersrepo "github.com/xantinium/gophermart/internal/repository/users"
)

type Options struct {
	UsersRepo  *usersrepo.UsersRepository
	TokensRepo *tokensrepo.TokensRepository
	OrdersRepo *ordersrepo.OrdersRepository
}

func New(opts Options) *UseCases {
	return &UseCases{
		usersRepo:  opts.UsersRepo,
		tokensRepo: opts.TokensRepo,
		ordersRepo: opts.OrdersRepo,
	}
}

type UseCases struct {
	usersRepo  *usersrepo.UsersRepository
	tokensRepo *tokensrepo.TokensRepository
	ordersRepo *ordersrepo.OrdersRepository
}
