package sectors

import (
	. "drift/common"
	"drift/ships"
	
	"fmt"

	. "github.com/klkblake/s3dm"
	
)


type ShipMap map[string]*ships.Ship

type Sector struct {
	Coords SectorCoords
	Name string
	ShipsByID ShipMap `msgpack:"-"`
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.Coords.X, sector.Coords.Y)
}

func SectorByCoords(x int, y int) *Sector {
	return &Sector{
		Coords: SectorCoords{X: x, Y: y},
		ShipsByID: make(ShipMap),
	}
}

func (sector *Sector) LoadShips(client StorageClient) {	
	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	for _, loc := range foundLocs {
		ship := loc.GetShip(client)
		sector.ShipsByID[ship.ID] = ship
		ship.Location.Body.Velocity = V3{0, 0, 0}
		ship.Location.Body.Spin = AxisAngle(V3{0, 0, 1}, 0.1)
	}
}

func (sector *Sector) Tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location.Body
		pos = pos.RK4Integrate(TICK_DELTA)
		ship.Location.Body = pos
	}
}

func prettyV3(vec V3) string {
	return fmt.Sprintf("<%.2f, %.2f, %.2f>", vec.X, vec.Y, vec.Z)
}

func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v:\n", sector.Name, sector.Coords)
	for _, ship := range sector.ShipsByID {
		fmt.Printf("   %s (%v, %v, %v, %v)\n", 
			ship.ID, 
			prettyV3(ship.Location.Body.Position),
			prettyV3(ship.Location.Body.Velocity),
			prettyV3(ship.Location.Body.Thrust),
			ship.Location.Body.Spin)
	}
}
