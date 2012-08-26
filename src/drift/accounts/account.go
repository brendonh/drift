package accounts

import (
	. "drift/common"
	"code.google.com/p/go.crypto/bcrypt"
)

type Account struct {
	Name string
	PasswordHash []byte
	Admin bool
}

// Storage API
func (account *Account) StorageKey() string {
	return account.Name
}

// User API
func (account *Account) DisplayName() string {
	return account.Name
}

// User API
func (account *Account) ID() string {
	return account.Name
}



func (account *Account) CheckPassword(given string) bool {
	var err = bcrypt.CompareHashAndPassword(
		account.PasswordHash,
		[]byte(given))
	return err == nil
}


func NewAccount(name string, password string) *Account {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &Account{name, hash, false}
}

// XXX BGH TODO: Serialize to avoid Riak races
func CreateAccount(name string, password string, context DriftServerContext) (*Account, bool) {
	var client = context.Storage()

	existing := &Account{Name: name}
	if client.Get(existing) {
		return nil, false
	}

	account := NewAccount(name, password)
	if !client.Put(account) {
		return nil, false
	}
	return account, true
}




