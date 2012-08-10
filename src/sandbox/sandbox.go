package main

import (
	"drift/storage"
	"drift/services"
	"drift/accounts"
	"drift/server"
	"drift/endpoints"

	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	flag.Parse()
	
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s [command]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	var command = flag.Arg(0)

	switch flag.Arg(0) {
	case "start":
		startServer()
	case "emptybucket":
		emptyBucket()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(2)
	}
}

func startServer() {
	fmt.Printf("Starting server...\n")
	var client = storage.NewRawRiakClient("http://localhost:8098")

	serviceCollection := services.NewServiceCollection()
	serviceCollection.AddService(accounts.GetService())

	var s = server.NewServer(client, serviceCollection)

	httpRpc := endpoints.NewHttpRpcEndpoint(":9999", s)
	s.AddEndpoint(httpRpc)		

	var stopper = make(chan os.Signal, 1)
	signal.Notify(stopper)

	s.Start()

	<-stopper
	close(stopper)

	fmt.Printf("Shutting down ...\n")
	s.Stop()
}

func emptyBucket() {
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s emptybucket [bucket]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	var client = storage.NewRawRiakClient("http://localhost:8098")
	var bucket = flag.Arg(1)

	fmt.Printf("Emptying bucket %s...\n", bucket)

	keys, ok := client.Keys(bucket)

	if !ok {
		fmt.Printf("Couldn't retrieve keys.\n", bucket)
		return
	}

	for _, key := range keys {
		fmt.Printf("Deleting key %s... ", key)
		ok := client.Delete(bucket, key)
		if ok {
			fmt.Printf("[ok]\n")
		} else {
			fmt.Printf("[error]\n")
		}
	}
}

	// sector := sectors.SectorByCoords(0, 0)
	// client.Get(sector)
	// fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	// sector.LoadShips(client)
	// sector.DumpShips()
	// sector.Tick()
	// sector.DumpShips()



