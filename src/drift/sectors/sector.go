package sectors

import (
	"drift/common"
	"drift/storage"
	"drift/ships"

	"fmt"
)


type ShipMap map[string]*ships.Ship

type Sector struct {
	Coords common.SectorCoords
	Name string
	ShipsByID ShipMap `msgpack:"-"`
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.Coords.X, sector.Coords.Y)
}

func SectorByCoords(x int, y int) *Sector {
	return &Sector{
		Coords: common.SectorCoords{X: x, Y: y},
		ShipsByID: make(ShipMap),
	}
}

func (sector *Sector) LoadShips(client storage.StorageClient) {	
	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	for _, loc := range foundLocs {
		ship := loc.GetShip(client)
		sector.ShipsByID[ship.ID] = ship
	}
}

func (sector *Sector) Tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location.Body
		pos = pos.RK4Integrate(common.TICK_DELTA)
		ship.Location.Body = pos
	}
}

func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v:\n", sector.Name, sector.Coords)
	for _, ship := range sector.ShipsByID {
		fmt.Printf("   %s (%v)\n", ship.ID, *ship.Location.Body)
	}
}