package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// ErrInvalidBody thrown when the body of a request can not be parsed
var ErrInvalidBody = errors.New("invalid body")

// MakeHTTPHandler returns all http handler for the Account service
func MakeHTTPHandler(logger log.Logger, endpoints Endpoints) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getAccountHandler := kithttp.NewServer(
		endpoints.GetByID,
		decodeGetAccountRequest,
		encodeResponse,
		options...,
	)

	getAccountsHandler := kithttp.NewServer(
		endpoints.GetList,
		decodeGetAccountsRequest,
		encodeResponse,
		options...,
	)

	updateAccountHandler := kithttp.NewServer(
		endpoints.Update,
		decodeUpdateAccountRequest,
		encodeResponse,
		options...,
	)

	createAccountHandler := kithttp.NewServer(
		endpoints.Create,
		decodeCreateAccountRequest,
		encodeCreateAccountResponse,
		options...,
	)

	deleteAccountHandler := kithttp.NewServer(
		endpoints.Delete,
		decodeDeleteAccountRequest,
		encodeDeleteAccountResponse,
		options...,
	)

	r := mux.NewRouter().PathPrefix("/accounts/").Subrouter().StrictSlash(true)

	r.Handle("/", getAccountsHandler).Methods("GET")
	r.Handle("/{id}", getAccountHandler).Methods("GET")
	r.Handle("/{id}", updateAccountHandler).Methods("PATCH")
	r.Handle("/", createAccountHandler).Methods("POST")
	r.Handle("/{id}", deleteAccountHandler).Methods("DELETE")

	return r
}

func decodeGetAccountsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	// TODO : parse pagination filter properly
	p := Pagination{Size: DefaultPaginationSize, Page: 0}

	f := Filter{
		IDs: strings.Split(r.URL.Query().Get("account_id"), ","),
	}

	return GetAccountsRequest{Filter: f, Pagination: p}, nil
}

func decodeGetAccountRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)

	return GetAccountRequest{ID: vars["id"]}, nil
}

func decodeUpdateAccountRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req UpdateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidBody
	}

	vars := mux.Vars(r)

	if len(req.Account.AccountID) == 0 {
		req.AccountID = vars["id"]
	}

	if req.Account.AccountID != vars["id"] {
		return nil, ErrInconsistentID
	}

	return req, nil
}

func decodeCreateAccountRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidBody
	}

	return req, nil
}

func decodeDeleteAccountRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)

	return DeleteAccountRequest{ID: vars["id"]}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeDeleteAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// TODO : refactor return
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeCreateAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// TODO : refactor return
	ID, ok := response.(string)
	if !ok {
		return errors.New("An error occured while creating Account")
	}
	w.Header().Set("Location", fmt.Sprintf("/accounts/%v", ID))
	w.WriteHeader(http.StatusCreated)
	return nil
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {

	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	case ErrInconsistentID,
		ErrInvalidBody:
		w.WriteHeader(http.StatusBadRequest)
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
