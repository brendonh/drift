package accounts

import (
	. "drift/common"
	"drift/services"

	"fmt"
)

// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetService() *services.Service {
	service := services.NewService("accounts")
	service.AddMethod(
		"register",
		[]APIArg{
		    APIArg{Name: "name", ArgType: StringArg},
		    APIArg{Name: "password", ArgType: StringArg},
	    },
		method_register)

	service.AddMethod(
		"login",
		[]APIArg{
		    APIArg{Name: "name", ArgType: StringArg},
		    APIArg{Name: "password", ArgType: StringArg},
	    },
		method_login)

	service.AddMethod(
		"ping",
		[]APIArg{},
		method_ping)
	
	return service
}


func method_register(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	account, ok := CreateAccount(
		args["name"].(string), 
		args["password"].(string), 
		context)
	
	if !ok {
		response["message"] = "User exists"
		return false, response
	}
	
	fmt.Printf(
		"Registered account: %s\n", 
		account.Name)
	
	return true, response
}


func method_login(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	session.Lock()
	defer session.Unlock()

	if session.User() != nil {
		response["message"] = "Already logged in"
		return false, response
	}

	var client = context.Storage()

	var account = &Account{Name: args["name"].(string)}

	if !client.Get(account) || !account.CheckPassword(args["password"].(string)) {
		response["message"] = "Invalid credentials"
		return false, response
	}

	session.SetUser(account)

	return true, response
}


func method_ping(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)
	response["message"] = "Pong"
	return true, response
}