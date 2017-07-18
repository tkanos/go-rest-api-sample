package account

import (
	"context"
	"errors"
)

// ErrNotFound is used when an Account is not found
var ErrNotFound = errors.New("Account not found")

// ErrInconsistentID ...
var ErrInconsistentID = errors.New("inconsistent Accountid")

// Service is the Order service interface
type Service interface {
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, filter Filter, pagination Pagination) ([]*Account, error)
	UpdateAccount(ctx context.Context, Account Account) error
	CreateAccount(ctx context.Context, Account Account) (string, error)
	DeleteAccount(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

// NewService return a new instance of order service
func NewService(r Repository) Service {
	return service{
		repository: r,
	}
}

// GetAccount returns an Account regarding the id passed in parameter
func (s service) GetAccount(ctx context.Context, id string) (a *Account, err error) {
	a, err = s.repository.GetAccount(ctx, id)

	if a == nil {
		err = ErrNotFound
	}

	return
}

// GetAccounts returns a list of Accounts regarding the ids passed in parameter
func (s service) GetAccounts(ctx context.Context, filter Filter, pagination Pagination) (accounts []*Account, err error) {

	accounts, err = s.repository.GetAccounts(ctx, filter, pagination)

	return
}

// UpdateAccount updates an existing Account
func (s service) UpdateAccount(ctx context.Context, a Account) error {
	return nil
}

// CreateAccount creates an Account
func (s service) CreateAccount(ctx context.Context, a Account) (id string, err error) {
	id, err = s.repository.CreateAccount(ctx, a)
	return
}

// DeleteAccount deletes an account
func (s service) DeleteAccount(ctx context.Context, id string) (err error) {
	err = s.repository.DeleteAccount(ctx, id)
	return
}
