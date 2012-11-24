package services

import (
	"drift/server"
	"drift/ships"

	"fmt"

	. "github.com/brendonh/go-service"
)

// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetSectorService() *Service {
	service := NewService("server")

	service.AddMethod(
		"control",
		[]APIArg{
  		  APIArg{Name: "id", ArgType: StringArg},
	    },
		method_control)

	return service
}


func method_control(args APIData, session Session, context ServerContext) (bool, APIData) {
	var server = context.(*server.DriftServer)
	var response = make(APIData)

	session.Lock()
	defer session.Unlock()

	var user = session.User()

	if user == nil {
		response["message"] = "Not logged in"
		return false, response
	}

	var loc = &ships.ShipLocation{ ShipID: args["id"].(string) }
	if !server.Storage().Get(loc) {
		response["message"] = "Ship not found"
		return false, response
	}

	sector, ok := server.SectorManager.Sectors[loc.Coords.String()]

	if !ok {
		fmt.Printf("No such sector\n")
		response["message"] = "Ship not in running sector"
		return false, response
	}

	ok, error := sector.Control(user, args["id"].(string), true)
	
	if !ok {
		response["message"] = error
		return false, response
	}

	//session.(DriftSession).SetAvatar(ship)

	return true, response
}