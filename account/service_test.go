package account

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewService_Should_Create_New_Service_Instance(t *testing.T) {
	fakeRepo := new(mockedAccountRepository)

	svc := NewService(fakeRepo)

	assert.NotNil(t, svc)
}

func Test_GetAccount_Should_Return_OK_If_Params_Is_Valid(t *testing.T) {
	expected := &Account{}
	AccountID := "12345"
	fakeRepo := new(mockedAccountRepository)
	fakeRepo.On("GetAccount", AccountID).Return(expected, nil)

	svc := NewService(fakeRepo)
	a, err := svc.GetAccount(context.Background(), AccountID)

	assert.NotNil(t, a)
	assert.Nil(t, err)
}

func Test_GetAccount_Should_Return_Error_If_Repository_Return_Error(t *testing.T) {
	expected := errors.New("error")
	AccountID := "12345"
	fakeRepo := new(mockedAccountRepository)
	fakeRepo.On("GetAccount", AccountID).Return(&Account{}, expected)

	svc := NewService(fakeRepo)
	_, err := svc.GetAccount(context.Background(), AccountID)

	assert.Equal(t, expected, err)
}

func Test_GetAccount_Should_Return_Error_If_Account_Is_Not_Found(t *testing.T) {
	expected := ErrNotFound
	AccountID := "12345"
	fakeRepo := new(mockedAccountRepository)
	fakeRepo.On("GetAccount", AccountID).Return(nil, nil)

	svc := NewService(fakeRepo)
	_, err := svc.GetAccount(context.Background(), AccountID)

	assert.Equal(t, expected, err)
}

func Test_GetAccounts_Should_Return_OK_If_Params_Is_Valid(t *testing.T) {
	expected := []*Account{}
	f := Filter{IDs: []string{"01234", "56789"}}
	p := Pagination{Size: 100}
	fakeRepo := new(mockedAccountRepository)
	fakeRepo.On("GetAccounts", f, p).Return(expected, nil)

	svc := NewService(fakeRepo)
	u, err := svc.GetAccounts(context.Background(), f, p)

	assert.NotNil(t, u)
	assert.Nil(t, err)
}
