package services

import (
	"drift/server"
	"drift/ships"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/brendonh/go-service"
	"github.com/brendonh/loge/src/loge"
)

// ------------------------------------------
// Service endpoints
// ------------------------------------------

func GetShipService() *Service {
	service := NewService("ships")

	service.AddMethod(
		"create",
		[]APIArg{
	      APIArg{Name: "name", ArgType: StringArg},
	    },
		method_create)

	service.AddMethod(
		"list",
		[]APIArg{},
		method_list)

	return service
}


func method_create(args APIData, session Session, context ServerContext) (bool, APIData) {
	var server = context.(*server.DriftServer)
	var response = make(APIData)

	var user = session.User()
	
	if user == nil {
		response["message"] = "Not logged in"
		return false, response
	}

	var db = server.DB()
	var id = uuid.New()
	ship := ships.NewShip(id, user.ID(), args["name"].(string))

	sector, ok := server.SectorManager.Ensure(0, 0)

	if !ok {
		response["message"] = "Home sector unavailable"
		return false, response
	}

	db.SetOne("ship", loge.LogeKey(id), ship)

	sector.Warp(ship, true)

	ship.SaveLocation(db)
	
	response["id"] = ship.ID
	return true, response
}


func method_list(args APIData, session Session, context ServerContext) (bool, APIData) {
	// var server = context.(DriftServerContext)
	var response = make(APIData)

	// var user = session.User()
	
	// if user == nil {
	// 	response["message"] = "Not logged in"
	// 	return false, response
	// }
	
	// var ship = &ships.Ship{ Owner: user.ID() }
	// var ships = make([]ships.Ship, 0)
	// //server.Storage().IndexLookup(ship, &ships, "Owner")

	// var shipInfo = make([]map[string]interface{}, len(ships))
	// for i, ship := range ships {
	// 	shipInfo[i] = make(map[string]interface{})
	// 	shipInfo[i]["id"] = ship.ID
	// 	shipInfo[i]["name"] = ship.Name
	// 	if ship.LoadLocation(server.Storage()) {
	// 		var sector = make(map[string]interface{})
	// 		sector["x"] = ship.Location.Coords.X
	// 		sector["y"] = ship.Location.Coords.Y
	// 		shipInfo[i]["sector"] = sector
	// 	}
	// }
	// response["ships"] = shipInfo

	return true, response
}
