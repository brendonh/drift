package main

import (
	"drift/storage"
	"drift/accounts"
	"drift/services"
	"drift/endpoints"
	"drift/ships"
	"drift/sectors"

	//"drift/simulation"
	//. "github.com/brendonh/s3dm-go"

	"fmt"
	"os"
	"os/signal"

)

func main() {
	var client = storage.NewRawRiakClient("http://localhost:8098")

	// ship, _ := ships.CreateShip("brendonh", "onemore", client)
	// fmt.Printf("Ship ID: %s\n", ship.ID)
	// return

	sector := sectors.SectorByCoords(0, 0)
	client.Get(sector)
	fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	for _, loc := range foundLocs {
		fmt.Printf("Loc: %s (%s)\n", loc.ShipID, loc.Body.Position)
		ship := loc.GetShip(client)
		fmt.Printf("Ship: %s (%v)\n", ship.Name, ship.Location)
	}

	// ship := &ships.Ship{ ID: "f974bd25-3349-4aff-9924-341171b5f2b3" }
	// client.Get(ship)

	// ship.LoadLocation(client)
	// fmt.Printf("Ship: %s\n", ship.Name)
	// fmt.Printf("Location: %v\n", ship.Location.Body.Position)

	// ship.SaveLocation(client)
	
	// p := &simulation.Powered{
	// 	Position: V3{0, 0, 0},
  	// 	Velocity: V3{0, 10, 0},
	// 	Thrust: V3{0, -10, 0},
	// }

	// for i := 0; i < 5; i++ {		
	// 	fmt.Printf("Tick %d: %v\n", i, p)
	// 	p = p.RK4Integrate(1.0)
	// }
}


	//var shipIDs []string
	//var ok bool

	// shipIDs, ok := client.Keys("Ship")

	// if ok {
	// 	for _, shipID := range shipIDs {
	// 		fmt.Printf("Deleting: %s\n", shipID)
	// 		client.Delete("Ship", shipID)
	// 	}
	// }

	// return


	//ships.CreateShip("brendonh", "onemore", client)

	// if ok {
	// 	fmt.Printf("Created ship: %s\n", ship.ID)
	// }

	// fmt.Printf("~~~~~~~~~~~~~~\n")
	
	// //searchShip := &ships.Ship{ Owner: "brendonh" }

	// foundShips := make([]ships.Ship, 5)

	// client.IndexLookup(searchShip, &foundShips, "Owner")	

	// for _, ship := range foundShips {
	// 	fmt.Printf("Ship: %s (%s)\n", ship.Name, ship.ID)
	// }


func startServer() {
	serviceCollection := services.NewServiceCollection()
	serviceCollection.AddService(accounts.GetService())

	var stopper = make(chan os.Signal, 1)
	signal.Notify(stopper)

	endpoint := endpoints.NewHttpRpcEndpoint(":9999", serviceCollection)

	fmt.Printf("Starting HTTP RPC: %v\n", endpoint.Start())

	<-stopper
	close(stopper)
	
	fmt.Printf("Shutting down ...\n")

	fmt.Printf("Stopping HTTP RPC: %v\n", endpoint.Stop())
}


// 	blob := bytes.NewBufferString(
// 		`{"service": "accounts", "method": "register", "data": {"email": "brendonh4@gmail.com", "password": "test"}}`).Bytes()	

// 	var args map[string]interface{}
// 	err := json.Unmarshal(blob, &args)
	
// 	if err != nil {
// 		fmt.Printf("Oh no: %s\n", err)
// 		return
// 	}

// 	var response = serviceCollection.HandleRequest(args)

// 	reply, _ := json.Marshal(response)

// 	fmt.Printf("Response: %s\n", reply)
// }


// func main() {
// 	client := drift.NewRiakClient("http://localhost:8098")

// 	sector := Sector{0, 1, "Away"}

// 	ok := client.Put(&sector)

// 	if !ok {
// 		fmt.Printf("Write Failed\n")
// 		return
// 	}

// 	fmt.Printf("Ok\n")

// 	newSector := Sector{X: 0, Y: 1}
// 	ok = client.Get(&newSector)

// 	if !ok {
// 		fmt.Printf("Read Failed\n")
// 		return
// 	}

// 	fmt.Printf("%s\n", newSector.Name)
// }
