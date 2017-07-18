package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeGetAccountEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("GetAccount", "1").Return(&(Account{}), nil)

	endpoint := MakeGetAccountEndpoint(fakeService)
	a, err := endpoint(nil, GetAccountRequest{ID: "1"})

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
	assert.NotNil(t, a)
}

func Test_MakeGetAccountsEndpoint(t *testing.T) {
	f := Filter{IDs: []string{"1", "2"}}
	p := Pagination{Size: 100}
	req := GetAccountsRequest{Filter: f, Pagination: p}

	fakeService := new(mockedService)
	fakeService.On("GetAccounts", f, p).Return([]*Account{&(Account{})}, nil)

	endpoint := MakeGetAccountsEndpoint(fakeService)
	a, err := endpoint(nil, req)

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
	assert.NotNil(t, a)
}

func Test_MakeUpdateAccountEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("UpdateAccount", Account{}).Return(nil, nil)

	endpoint := MakeUpdateAccountEndpoint(fakeService)
	_, err := endpoint(nil, UpdateAccountRequest{})

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
}

func Test_MakeCreateAccountEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("CreateAccount", Account{}).Return("", nil)

	endpoint := MakeCreateAccountEndpoint(fakeService)
	_, err := endpoint(nil, CreateAccountRequest{})

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
}
func Test_MakeDeleteAccountEndpoint(t *testing.T) {
	fakeService := new(mockedService)
	fakeService.On("DeleteAccount", "1").Return(nil)

	endpoint := MakeDeleteAccountEndpoint(fakeService)
	_, err := endpoint(nil, DeleteAccountRequest{ID: "1"})

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
}
