package account

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_MakeHTTPHandler(t *testing.T) {
	h := MakeHTTPHandler(log.NewNopLogger(), Endpoints{})

	assert.NotNil(t, h)
}

func Test_DecodeGetAccountRequest(t *testing.T) {
	expected := GetAccountRequest{}
	r, _ := http.NewRequest("GET", "/Accounts/1", nil)

	req, err := decodeGetAccountRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}

func Test_DecodeGetAccountsRequest(t *testing.T) {
	f := Filter{IDs: []string{"1", "2"}}
	p := Pagination{Size: 100}
	expected := GetAccountsRequest{Filter: f, Pagination: p}
	r, _ := http.NewRequest("GET", "/Accounts/?id=1,2", nil)

	req, err := decodeGetAccountsRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}

func Test_DecodeUpdateAccountRequest(t *testing.T) {
	expected := UpdateAccountRequest{}
	r, _ := http.NewRequest("PUT", "/Accounts/1", bytes.NewBufferString("{}"))

	req, err := decodeUpdateAccountRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}
func Test_DecodeUpdateAccountRequest_Should_Returns_ErrInvalidAccountID_When_AccountID_In_Body_Not_Match(t *testing.T) {
	expected := ErrInconsistentID
	r, _ := http.NewRequest("PUT", "/Accounts/1", bytes.NewBufferString("{\"Account_id\":\"2\"}\n"))

	req, err := decodeUpdateAccountRequest(context.Background(), r)

	assert.Nil(t, req)
	assert.Equal(t, expected, err)
}

func Test_DecodeUpdateAccountRequest_Should_Returns_ErrInvalidBody_When_Body_Is_Invalid(t *testing.T) {
	expected := ErrInvalidBody
	r, _ := http.NewRequest("PUT", "/Accounts/1", bytes.NewBufferString("invalidjson"))

	_, err := decodeUpdateAccountRequest(context.Background(), r)

	assert.Equal(t, expected, err)
}

func Test_DecodeCreateAccountRequest(t *testing.T) {
	expected := CreateAccountRequest{}
	r, _ := http.NewRequest("POST", "/Accounts/", bytes.NewBufferString("{}"))

	req, err := decodeCreateAccountRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}

func Test_DecodeCreateAccountRequest_Should_Returns_ErrInvalidBody_When_Body_Is_Invalid(t *testing.T) {
	expected := ErrInvalidBody
	r, _ := http.NewRequest("POST", "/Accounts/", bytes.NewBufferString("invalidjson"))

	_, err := decodeCreateAccountRequest(context.Background(), r)

	assert.Equal(t, expected, err)
}

func Test_DecodeDeleteAccountRequest(t *testing.T) {
	expected := DeleteAccountRequest{}
	r, _ := http.NewRequest("DELETE", "/Accounts/1", nil)

	req, err := decodeDeleteAccountRequest(context.Background(), r)

	assert.Nil(t, err)
	assert.Equal(t, expected, req)
}

func Test_EncodeResponse(t *testing.T) {
	response := struct {
		ID string
	}{"12345"}
	expected := "{\"ID\":\"12345\"}\n"

	w := httptest.NewRecorder()
	err := encodeResponse(context.Background(), w, response)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(body))
}

func Test_EncodeResponse_Should_Return_JSON_ContentType(t *testing.T) {
	expected := "application/json; charset=utf-8"

	w := httptest.NewRecorder()
	err := encodeResponse(context.Background(), w, nil)

	assert.Nil(t, err)
	assert.Equal(t, expected, w.Header().Get("Content-Type"))
}

func Test_EncodeCreateAccountResponse(t *testing.T) {
	ID := "123"
	response := ID
	expectedBody := ""
	expectedLocation := fmt.Sprintf("/accounts/%v", ID)

	w := httptest.NewRecorder()
	err := encodeCreateAccountResponse(context.Background(), w, response)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))
	assert.Equal(t, expectedLocation, w.Header().Get("Location"))
	assert.Equal(t, http.StatusCreated, w.Code)
}

func Test_EncodeDeleteAccountResponse(t *testing.T) {

	w := httptest.NewRecorder()
	err := encodeDeleteAccountResponse(context.Background(), w, nil)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_EncodeError_Should_Return_JSON_ContentType(t *testing.T) {
	expected := "application/json; charset=utf-8"

	w := httptest.NewRecorder()
	encodeError(context.Background(), errors.New("error"), w)

	assert.Equal(t, expected, w.Header().Get("Content-Type"))
}

func Test_EncodeError(t *testing.T) {
	err := errors.New("fake error")
	expected := "{\"error\":\"fake error\"}\n"

	w := httptest.NewRecorder()
	encodeError(context.Background(), err, w)
	body, err := ioutil.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(body))
}

func Test_EncodeError_Should_Correctly_Map_Error(t *testing.T) {
	var flagtests = []struct {
		in  error
		out int
	}{
		{errors.New("not handled error"), http.StatusInternalServerError},
		{ErrInvalidBody, http.StatusBadRequest},
		{ErrInconsistentID, http.StatusBadRequest},
		{ErrNotFound, http.StatusNotFound},
	}

	for _, tt := range flagtests {
		w := httptest.NewRecorder()
		encodeError(context.Background(), tt.in, w)

		assert.Equal(t, tt.out, w.Code)
	}
}
