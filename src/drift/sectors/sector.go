package sectors

import (
	. "drift/common"
	"drift/ships"
	"drift/simulation"
	
	"fmt"
	"time"
)


type ShipMap map[string]*ships.Ship

type Sector struct {
	Coords SectorCoords
	Name string

	ShipsByID ShipMap         `msgpack:"-"`
	manager *SectorManager    `msgpack:"-"`
	chanStop chan int         `msgpack:"-"`
	chanTick <-chan time.Time `msgpack:"-"`
	bodies [1000]simulation.PoweredBody
}

func (sector *Sector) StorageKey() string {
	return sector.Coords.String()
}

func SectorByCoords(x int, y int, manager *SectorManager) *Sector {
	return &Sector{
		Coords: SectorCoords{X: x, Y: y},
		ShipsByID: make(ShipMap),

		manager: manager,
		chanStop: make(chan int, 0),
	}
}

func (sector *Sector) Start() {
	sector.loadShips()
	//sector.DumpShips()
	sector.chanTick = time.Tick(time.Duration(TICK_DELTA) * time.Millisecond)
	go sector.loop()
}

func (sector *Sector) Stop() {
	sector.chanStop <- 1
	<- sector.chanStop
}

func (sector *Sector) loop() {
	fmt.Printf("Sector started: %s\n", sector.Coords.String())

	for {
		select {
		case <-sector.chanStop:
			fmt.Printf("Sector stopping %s\n", sector.Coords.String())
			sector.chanStop <- 1
			break

		case <-sector.chanTick:
			//var start = time.Now()
			sector.tick()
			//fmt.Printf("Tick: %v\n", time.Since(start))
		}
	}
}


func (sector *Sector) loadShips() {	
	var client = sector.manager.context.Storage()

	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	var chunkSize = 32
	var totalShips = len(foundLocs)
	fmt.Printf("Loading %d ships...\n", totalShips)
	var start = time.Now()

	var chunkAcks = make(chan int)
	var shipStream = make(chan *ships.Ship)

	var chunks = 0

	for chunk := 0; chunk <= totalShips / chunkSize; chunk++ {
		var start = chunk * chunkSize
		var end = (chunk + 1) * chunkSize

		if end > totalShips {
			end = totalShips
		}

		if end <= start {
			break
		}

		go sector.loadShipChunk(foundLocs[start:end], client, chunkAcks, shipStream)
		chunks++;
	}

	var finishedChunks = 0
	var shipsLoaded = 0
	var ship *ships.Ship

	for finishedChunks < chunks {
		select {
		case <-chunkAcks:
			finishedChunks += 1
		case ship = <-shipStream:
			shipsLoaded += 1
			sector.ShipsByID[ship.ID] = ship
		}
	}

	fmt.Printf("Loaded %d ships in %v\n", shipsLoaded, time.Since(start))
}

func (sector *Sector) loadShipChunk(
	locs []ships.ShipLocation, client StorageClient, 
	done chan int, stream chan *ships.Ship) {
	for i := range locs {
		var loc = locs[i]
		stream <- loc.GetShip(client)
	}
	done <- 1
}


func (sector *Sector) tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location.Body
		pos = pos.EulerIntegrate(1.0)
		ship.Location.Body = pos
	}
}


func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v (%d):\n", sector.Name, sector.Coords, len(sector.ShipsByID))
	for _, ship := range sector.ShipsByID {
		ship.Dump()
	}
}
