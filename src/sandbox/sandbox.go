package main

import (
	"drift/storage"
	"drift/services"
	"drift/accounts"
	"drift/endpoints"
	_ "drift/sectors"
	"drift/ships"
	"drift/server"
	"drift/simulation"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"math/rand"
	"github.com/klkblake/s3dm"
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
	case "sandbox":
		sandbox()
	case "createShips":
		createShips()
	case "emptyBucket":
		emptyBucket()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(2)
	}
}

func startServer() {
	fmt.Printf("Starting server...\n")

	runtime.GOMAXPROCS(4)

	var s = buildServer()

	s.AddEndpoint(endpoints.NewHttpRpcEndpoint(":9999", s))
	s.AddEndpoint(endpoints.NewWebsocketEndpoint(":9998", s))

	var stopper = make(chan os.Signal, 1)
	signal.Notify(stopper)

	s.Start()

	s.SectorManager.Ensure(999, 999)

	<-stopper
	close(stopper)

	fmt.Printf("Shutting down ...\n")
	s.Stop()
}


func buildServer() *server.Server {
	var client = storage.NewRawRiakClient("http://localhost:8098")

	serviceCollection := services.NewServiceCollection()
	serviceCollection.AddService(accounts.GetService())
	serviceCollection.AddService(ships.GetService())

	return server.NewServer(client, serviceCollection)
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

func sandbox() {
	var server = buildServer()
	var manager = server.SectorManager

	manager.Ensure(999, 999)

	// client.Get(sector)
	// fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	// sector.LoadShips(client)

	// sector.DumpShips()
}


func createShips() {
	var server = buildServer()
	var client = server.Storage()

	var manager = server.SectorManager

	var sector, ok = manager.Ensure(999, 999)

	if !ok {
		fmt.Printf("Sector load error\n")
		return
	}

	fmt.Printf("Sector: %v\n", sector)

	var account = &accounts.Account{Name: "sandbox"}
	ok = client.Get(account)

	if !ok {
		fmt.Printf("User load error\n")
		return
	}

	fmt.Printf("User: %v\n", account.Name)

	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("sandbox%03d", i)
		var id = client.GenerateID()

		var location = &ships.ShipLocation {
			ShipID: id,
			Coords: sector.Coords,
			Body: new(simulation.PoweredBody),
		}

		var body = location.Body
		body.Position.X = float64(rand.Intn(1000))
		body.Position.Y = float64(rand.Intn(1000))
		body.Velocity.X = rand.Float64()
		body.Velocity.Y = rand.Float64()
		body.Spin = s3dm.AxisAngle(s3dm.V3{0, 0, 1}, rand.Float64())

		ship := ships.NewShip(id, account.ID(), name)
		ship.Location = location
		fmt.Printf("%v || %v || %v\n", ship, ship.Location, ship.Location.Body)
		client.Put(ship)
		ship.SaveLocation(client)
	}
	
	// client.Get(sector)
	// fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	// sector.LoadShips(client)

	// sector.DumpShips()

	// for i := 0; i < 100; i++ {
	// 	sector.Tick()
	// 	sector.DumpShips()
	// }

}