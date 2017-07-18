package account

import (
	"context"
)

// Repository represents an user repository interface
type Repository interface {
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, filter Filter, pagination Pagination) ([]*Account, error)
	UpdateAccount(ctx context.Context, u Account) error
	CreateAccount(ctx context.Context, u Account) (string, error)
	DeleteAccount(ctx context.Context, id string) error
}
