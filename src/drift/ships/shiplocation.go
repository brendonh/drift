package ships

import (
	. "drift/common"
	"drift/simulation"
)

type ShipLocation struct {
	ShipID string
	Coords SectorCoords `indexed:"true"`
	Body *simulation.PoweredBody
}

func (loc *ShipLocation) StorageKey() string {
	return loc.ShipID
}

func (loc *ShipLocation) GetShip(client StorageClient) *Ship {
	ship := &Ship{ ID: loc.ShipID }
	ok := client.Get(ship)
	if !ok {
		return nil
	}
	ship.Location = loc
	return ship
}