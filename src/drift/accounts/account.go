package accounts

import (
	"drift/storage"
	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Name string
	PasswordHash []byte
	Admin bool
}

func NewUser(name string, password string) *User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{name, hash, false}
}

// XXX BGH TODO: Serialize to avoid Riak races
func CreateUser(name string, password string, client storage.StorageClient) (*User, bool) {
	user := NewUser(name, password)
	if !client.Put(user) {
		return nil, false
	}
	return user, true
}

func (user *User) StorageKey() string {
	return user.Name
}

