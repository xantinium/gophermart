package usecases

import (
	"github.com/xantinium/gophermart/internal/repository/tokensrepo"
	usersrepo "github.com/xantinium/gophermart/internal/repository/users"
)

type UseCasesOptions struct {
	UsersRepo  *usersrepo.UsersRepository
	TokensRepo *tokensrepo.TokensRepository
}

func NewUseCases(opts UseCasesOptions) *UseCases {
	return &UseCases{
		usersRepo:  opts.UsersRepo,
		tokensRepo: opts.TokensRepo,
	}
}

type UseCases struct {
	usersRepo  *usersrepo.UsersRepository
	tokensRepo *tokensrepo.TokensRepository
}
