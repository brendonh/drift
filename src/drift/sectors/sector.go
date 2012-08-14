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
			var start = time.Now()
			sector.tick()
			fmt.Printf("Tick: %v\n", time.Since(start))
			sector.DumpShips()
		}
	}
}


func (sector *Sector) loadShips() {	
	var client = sector.manager.context.Storage()

	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	fmt.Printf("Found %d ships\n", len(foundLocs))

	for i := range foundLocs {
		var loc = foundLocs[i]
		ship := loc.GetShip(client)
		fmt.Printf("Loaded ship '%s'\n", ship.Name)
		sector.ShipsByID[ship.ID] = ship

		// sector.bodies[i] = *ship.Location.Body
		// ship.Location.Body = &sector.bodies[i]

	}

	sector.DumpShips()
}

func (sector *Sector) tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location.Body
		pos = pos.RK4Integrate(1.0)
		ship.Location.Body = pos
	}
	// var i1 = make(chan int)
	// var i2 = make(chan int)
	//sector.tickRange(0, 1000)
	// go sector.tickRange(500, 1000, i2)
	// <-i1
	// <-i2
}

func (sector *Sector) tickRange(start int, stop int) {
	for i := start; i < stop; i++ {
		sector.bodies[i] = *sector.bodies[i].RK4Integrate(1.0)
	}
}


func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v:\n", sector.Name, sector.Coords)
	for _, ship := range sector.ShipsByID {
		fmt.Printf("   %s (%v, %v, %v, %v)\n", 
			ship.ID, 
			ship.Location.Body.Position.String(),
			ship.Location.Body.Velocity.String(),
			ship.Location.Body.Thrust.String(),
			ship.Location.Body.Spin.String())
	}
}
