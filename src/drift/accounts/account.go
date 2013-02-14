package accounts

import (
	. "drift/common"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/brendonh/loge/src/loge"
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

func CreateAccount(name string, password string, context DriftServerContext) (*Account, bool) {
	var db = context.DB()
	var account *Account = NewAccount(name, password)
	var success bool

	db.Transact(func (t *loge.Transaction) {
		if !t.Exists("account", loge.LogeKey(name)) {
			t.Set("account", loge.LogeKey(name), account)
			success = true
		}
	}, 0)

	if success {
		return account, true
	} 

	return nil, false
}




