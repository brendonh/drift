package services

import (
	"drift/server"
	"drift/endpoints"
	"drift/control"
	"drift/simulation"

	"fmt"

	. "github.com/brendonh/go-service"
	"github.com/brendonh/loge/src/loge"
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

	var db = server.DB()
	var loc = db.ReadOne("shiplocation", args["id"].(loge.LogeKey)).(*simulation.PoweredBody)
	if loc == nil {
		response["message"] = "Ship not found"
		return false, response
	}

	sector, ok := server.SectorManager.Sectors[loc.Coords.String()]

	if !ok {
		fmt.Printf("No such sector\n")
		response["message"] = "Ship not in running sector"
		return false, response
	}

	ok, error := sector.Control(
		session.(*endpoints.ServerSession), 
		args["id"].(string), 
		true, control.ControlSpec(true))

	if !ok {
		response["message"] = error
		return false, response
	}

	//session.(DriftSession).SetAvatar(ship)

	//session.Send([]byte("Hello World"))

	return true, response
}