package mongoDb

import (
	"context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tkanos/go-rest-api-sample/account"
)

const emailExpirationHours = 24

type accountRepository struct {
	session *mgo.Session
}

// NewAccountRepository creates a new instance of a legacy account repository
func NewAccountRepository(s *mgo.Session) (account.Repository, error) {

	err := ensureIndex(s)

	return accountRepository{
		session: s,
	}, err
}

func ensureIndex(s *mgo.Session) error {
	session := s.Copy()
	defer session.Close()

	c := session.DB("store").C("accounts")

	index := mgo.Index{
		Key:        []string{"account_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return c.EnsureIndex(index)
}

// Getaccount ...
func (r accountRepository) GetAccount(ctx context.Context, id string) (a *account.Account, err error) {
	session := r.session.Copy()
	defer session.Close()

	c := session.DB("store").C("accounts")

	err = c.Find(bson.M{"account_id": id}).One(&a)
	return
}

// Getaccounts ...
func (r accountRepository) GetAccounts(ctx context.Context, filter account.Filter, pagination account.Pagination) (accounts []*account.Account, err error) {
	session := r.session.Copy()
	defer session.Close()

	c := session.DB("store").C("accounts")

	m := bson.M{}

	if len(filter.IDs) > 0 {
		m = bson.M{"account_id": bson.M{"$in": filter.IDs}}
	}

	err = c.Find(m).Skip(pagination.Page).Limit(pagination.Size).All(&accounts)

	return
}

// Updateaccount ...
func (r accountRepository) UpdateAccount(ctx context.Context, a account.Account) error {
	return nil
}

// Createaccount ...
func (r accountRepository) CreateAccount(ctx context.Context, a account.Account) (string, error) {
	session := r.session.Copy()
	defer session.Close()

	a.AccountID = bson.NewObjectId().Hex()
	c := session.DB("store").C("accounts")

	err := c.Insert(a)

	return a.AccountID, err
}

// Deleteaccount ...
func (r accountRepository) DeleteAccount(ctx context.Context, id string) error {
	session := r.session.Copy()
	defer session.Close()

	c := session.DB("store").C("accounts")

	return c.Remove(bson.M{"account_id": id})
}
