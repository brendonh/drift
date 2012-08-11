package ships

import (
	. "drift/common"
	"drift/services"
)

// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetService() *services.Service {
	service := services.NewService("ships")
	service.AddMethod(
		"create",
		[]APIArg{
	      APIArg{Name: "name", ArgType: StringArg},
	    },
		method_create)

	service.AddMethod(
		"list",
		[]APIArg{},
		method_register)

	service.AddMethod(
		"control",
		[]APIArg{
		  APIArg{Name: "id", ArgType: StringArg},
	    },
		method_control)

	return service
}


func method_create(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	var user = session.User()
	
	if user == nil {
		response["message"] = "Not logged in"
		return false, response
	}
	
	ship, ok := CreateShip(user.ID(), args["name"].(string), context)

	if !ok {
		return false, response
	}

	response["id"] = ship.ID
	return true, response
}


func method_register(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	var user = session.User()
	
	if user == nil {
		response["message"] = "Not logged in"
		return false, response
	}
	
	var ship = &Ship{ Owner: user.ID() }
	var ships = make([]Ship, 0)
	context.Storage().IndexLookup(ship, &ships, "Owner")

	var shipInfo = make([]map[string]interface{}, len(ships))
	for i, ship := range ships {
		shipInfo[i] = make(map[string]interface{})
		shipInfo[i]["id"] = ship.ID
		shipInfo[i]["name"] = ship.Name
	}
	response["ships"] = shipInfo

	return true, response
}


func method_control(args APIData, session Session, context ServerContext) (bool, APIData) {
	var response = make(APIData)

	session.Lock()
	defer session.Unlock()

	var user = session.User()
	
	if user == nil {
		response["message"] = "Not logged in"
		return false, response
	}

	var ship = &Ship{ ID: args["id"].(string) }
	var ok = context.Storage().Get(ship)
	
	if !ok {
		response["message"] = "No such ship"
		return false, response
	}

	if ship.Owner != user.ID() {
		response["message"] = "Not your ship"
		return false, response
	}

	session.SetAvatar(ship)

	return true, response
}