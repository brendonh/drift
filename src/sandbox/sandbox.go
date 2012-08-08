package main

import (
	"drift/storage"
	"drift/accounts"
	"drift/services"
	"drift/endpoints"
	"drift/sectors"

	"fmt"
	"os"
	"os/signal"

)

func main() {
	var client = storage.NewRawRiakClient("http://localhost:8098")

	sector := sectors.SectorByCoords(0, 0)
	client.Get(sector)
	fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	sector.LoadShips(client)
	sector.DumpShips()
	sector.Tick()
	sector.DumpShips()
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

