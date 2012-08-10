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
	
	return service
}


func method_register(args APIData, context ServerContext) (bool, APIData) {
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


func method_login(args APIData, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	var client = context.Storage()

	var account = Account{Name: args["name"].(string)}

	if !client.Get(&account) || !account.CheckPassword(args["password"].(string)) {
		response["message"] = "Invalid credentials"
		return false, response
	}

	return true, response
}