package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints represent the order service endpoints
type Endpoints struct {
	GetByID endpoint.Endpoint
	GetList endpoint.Endpoint
	Update  endpoint.Endpoint
	Create  endpoint.Endpoint
	Delete  endpoint.Endpoint
}

// MakeGetAccountEndpoint returns an endpoint used for getting an account
func MakeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountRequest)

		return s.GetAccount(ctx, req.ID)
	}
}

// MakeGetAccountsEndpoint returns an endpoint used for getting accounts
func MakeGetAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountsRequest)

		return s.GetAccounts(ctx, req.Filter, req.Pagination)
	}
}

// MakeUpdateAccountEndpoint returns an endpoint used for updating an account
func MakeUpdateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateAccountRequest)

		return nil, s.UpdateAccount(ctx, req.Account)
	}
}

// MakeCreateAccountEndpoint returns an endpoint used for creating an account
func MakeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)

		return s.CreateAccount(ctx, req.Account)
	}
}

// MakeDeleteAccountEndpoint returns an endpoint used for deleting an account
func MakeDeleteAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteAccountRequest)

		return nil, s.DeleteAccount(ctx, req.ID)
	}
}

// GetAccountRequest represents the request parameters used for getting one Account
type GetAccountRequest struct {
	ID string `json:"id"`
}

// GetAccountsRequest represents the request parameters used for getting Accounts
type GetAccountsRequest struct {
	Filter
	Pagination
}

// UpdateAccountRequest represents the request parameters used for updating Account
type UpdateAccountRequest struct {
	Account
}

// CreateAccountRequest represents the request parameters used for creating Account
type CreateAccountRequest struct {
	Account
}

// DeleteAccountRequest represents the request parameters used for delete an account
type DeleteAccountRequest struct {
	ID string `json:"id"`
}

// Pagination ...
type Pagination struct {
	Size int
	Page int
}

// DefaultPaginationSize ...
const DefaultPaginationSize int = 100

// Filter ...
type Filter struct {
	IDs []string
}
