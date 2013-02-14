package main

import (
	"drift/accounts"
	"drift/endpoints"
	"drift/ships"
	"drift/server"
	"drift/simulation"
	"drift/services"
	"drift/sectors"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"math"
	"math/rand"

	"code.google.com/p/go-uuid/uuid"
	"github.com/brendonh/go-service"
	"github.com/brendonh/loge/src/loge"
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

	var manager = s.SectorManager
	var _, ok = manager.Ensure(0, 0)
	if !ok {
		_, ok = manager.Create(0, 0, "Home")
		if !ok {
			fmt.Printf("Sector creation error\n")
			return
		}
	}

	s.AddEndpoint(goservice.NewHttpRpcEndpoint(":9999", s, nil))

	var websocketEndpoint = goservice.NewWebsocketEndpoint(":9998", s)
	websocketEndpoint.Handler = endpoints.DriftMessageHandler
	s.AddEndpoint(websocketEndpoint)

	s.AddEndpoint(goservice.NewTelnetEndpoint(":6060", s))

	var stopper = make(chan os.Signal, 1)
	signal.Notify(stopper)

	s.Start()
	defer s.Stop()

	_, ok = s.SectorManager.Ensure(0, 0)
	if !ok {
		return
	}

	<-stopper
	close(stopper)

	fmt.Printf("Shutting down ...\n")
	s.DB().Close()
}


func buildServer() *server.DriftServer {
	var db = loge.NewLogeDB(loge.NewLevelDBStore("data/sandbox"))

	db.CreateType(loge.NewTypeDef("sector", 1, &sectors.Sector{}))
	db.CreateType(loge.NewTypeDef("account", 1, &accounts.Account{}))
	db.CreateType(loge.NewTypeDef("ship", 1, &ships.Ship{}))

	var locDef = loge.NewTypeDef("shiplocation", 1, &simulation.PoweredBody{})
	locDef.Links = loge.LinkSpec{ "sector": "sector" }
	db.CreateType(locDef)

	serviceCollection := goservice.NewServiceCollection()
	serviceCollection.AddService(services.GetAccountService())
	serviceCollection.AddService(services.GetSectorService())
	serviceCollection.AddService(services.GetShipService())
	serviceCollection.AddService(loge.GetService())

	return server.NewDriftServer(db, serviceCollection)
}


func emptyBucket() {
	// if flag.NArg() != 2 {
	// 	fmt.Fprintf(os.Stderr, "usage: %s emptybucket [bucket]\n", os.Args[0])
	// 	flag.PrintDefaults()
	// 	os.Exit(2)
	// }

	// var client = storage.NewRawRiakClient("http://localhost:8098")
	// var bucket = flag.Arg(1)

	// fmt.Printf("Emptying bucket %s...\n", bucket)

	// keys, ok := client.Keys(bucket)

	// if !ok {
	// 	fmt.Printf("Couldn't retrieve keys.\n", bucket)
	// 	return
	// }

	// for _, key := range keys {
	// 	fmt.Printf("Deleting key %s... ", key)
	// 	ok := client.Delete(bucket, key)
	// 	if ok {
	// 		fmt.Printf("[ok]\n")
	// 	} else {
	// 		fmt.Printf("[error]\n")
	// 	}
	// }
}


func sandbox() {
	var state = simulation.NewSectorState()

	fmt.Printf("Bodies: %d\n", len(state.Bodies))
	for i := 0; i < 1000; i++ {
		state = state.AddBody(simulation.PoweredBody{ ShipID: fmt.Sprintf("foo%d", i) })
	}
	fmt.Printf("Bodies: %d\n", len(state.Bodies))

	var start = time.Now()
	state = state.RemoveBody("foo850")
	fmt.Printf("Remove: %v\n", time.Since(start))

	fmt.Printf("Bodies: %d\n", len(state.Bodies))
}


func createShips() {
	var server = buildServer()

	var manager = server.SectorManager

	var sector, ok = manager.Ensure(0, 0)

	if !ok {
		sector, ok = manager.Create(0, 0, "Home")
		if !ok {
			fmt.Printf("Sector creation error\n")
			return
		}
	}

	fmt.Printf("Sector: %v\n", sector.Name)

	var db = server.DB()
	var account = db.ReadOne("account", "sandbox").(*accounts.Account)

	if account == nil {
		account, ok = accounts.CreateAccount("sandbox", "password", server)
		if !ok {
			fmt.Printf("User load error\n")
			return
		}
	}

	fmt.Printf("User: %v\n", account.Name)

	var sectorLink = []loge.LogeKey{ loge.LogeKey(sector.StorageKey()) }
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("sandbox%03d", i)
		var id = uuid.New()

		var body = &simulation.PoweredBody {
			ShipID: id,
			Coords: sector.Coords,
		}

		body.Position.X = float64(rand.Intn(1000))
		body.Position.Y = float64(rand.Intn(1000))
		body.Velocity.X = rand.Float64()
		body.Velocity.Y = rand.Float64()
		body.Thrust.X = rand.Float64()
		body.Thrust.Y = rand.Float64()

		var rot = math.Pi / 20
		body.Spin.X = math.Cos(rot)
		body.Spin.Y = math.Sin(rot)

		ship := ships.NewShip(id, account.ID(), name)
		ship.Location = body

		db.Transact(func (t *loge.Transaction) {
			var id = loge.LogeKey(id)
			t.Set("ship", id, ship)
			t.Set("shiplocation", id, ship.Location)
			t.SetLinks("shiplocation", "sector", id, sectorLink)
		}, 0)
	}

	fmt.Printf("Sector: %s (%d, %d)\n", sector.Name, sector.Coords.X, sector.Coords.Y)

	sector.DumpShips()

	// for i := 0; i < 100; i++ {
	// 	sector.Tick()
	// 	sector.DumpShips()
	// }

}