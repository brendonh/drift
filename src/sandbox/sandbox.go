package main

import (
	"drift/storage"
	"drift/accounts"
	"drift/services"
	"drift/endpoints"
	"drift/ships"

	"fmt"
	"os"
	"os/signal"

)

type Sector struct {
	X, Y int
	Name string
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.X, sector.Y)
}

func main() {
	client := storage.NewRiakClient("http://localhost:8098")	
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
	
	searchShip := &ships.Ship{ Owner: "brendonh" }

	foundShips := make([]ships.Ship, 0)

	client.IndexLookup(searchShip, &foundShips, "Owner")	

	for _, ship := range foundShips {
		fmt.Printf("Ship: %s (%s)\n", ship.Name, ship.ID)
	}

}

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
