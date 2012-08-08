package accounts

import (
	. "drift/common"
	"drift/services"

	"fmt"

	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Name string
	PasswordHash []byte
	Admin bool
}

func (user *User) StorageKey() string {
	return user.Name
}


func NewUser(name string, password string) *User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{name, hash, false}
}

// XXX BGH TODO: Serialize to avoid Riak races
func CreateUser(name string, password string, context *ServerContext) (*User, bool) {
	var client = context.StorageClient

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


func (user *User) CheckPassword(given string) bool {
	var err = bcrypt.CompareHashAndPassword(
		user.PasswordHash,
		[]byte(given))
	return err == nil
}

// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetService() *services.Service {
	service := services.NewService("accounts")
	service.AddMethod(
		"register",
		[]APIArg{
		    APIArg{Name: "name", ArgType: services.StringArg},
		    APIArg{Name: "password", ArgType: services.StringArg},
	    },
		method_register)

	service.AddMethod(
		"login",
		[]APIArg{
		    APIArg{Name: "name", ArgType: services.StringArg},
		    APIArg{Name: "password", ArgType: services.StringArg},
	    },
		method_login)
	
	return service
}


func method_register(args APIData, context *ServerContext) (bool, APIData) {
	var response = make(APIData)

	user, ok := CreateUser(
		args["name"].(string), 
		args["password"].(string), 
		context)
	
	if !ok {
		response["message"] = "User exists"
		return false, response
	}
	
	fmt.Printf(
		"Registered account: %s\n", 
		user.Name)
	
	return true, response
}

func method_login(args APIData, context *ServerContext) (bool, APIData) {
	var response = make(APIData)

	var client = context.StorageClient

	var user = User{Name: args["name"].(string)}

	if !client.Get(&user) || !user.CheckPassword(args["password"].(string)) {
		response["message"] = "Invalid credentials"
		return false, response
	}

	return true, response
}