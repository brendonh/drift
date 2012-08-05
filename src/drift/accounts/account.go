package accounts

import (
	"drift/storage"
	"drift/services"

	"fmt"

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
	existing := User{Name: name}
	if client.Get(&existing) {
		return nil, false
	}

	user := NewUser(name, password)
	if !client.Put(user) {
		return nil, false
	}
	return user, true
}

func (user *User) StorageKey() string {
	return user.Name
}


// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetService() *services.Service {
	service := services.NewService("accounts")
	service.AddMethod(
		"register",
		[]services.Arg{
		    services.Arg{Name: "email", ArgType: services.StringArg},
		    services.Arg{Name: "password", ArgType: services.StringArg},
	    },
		method_register)
	
	return service
}


func method_register(args map[string]interface{}) (bool, map[string]interface{}) {
	var response = make(map[string]interface{})

	// XXX BGH TODO: Get this somehow
	var client = storage.NewRiakClient("http://localhost:8098")

	user, ok := CreateUser(
		args["email"].(string), 
		args["password"].(string), 
		client)
	
	if !ok {
		response["message"] = "User exists"
		return false, response
	}
	
	fmt.Printf(
		"Registered account: %s\n", 
		user.Name)
	
	return true, response
}
