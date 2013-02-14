package sectors

import (
	. "drift/common"
	"drift/ships"
	"drift/simulation"
	"drift/endpoints"
	"drift/control"
	
	"fmt"
	"time"
)


type ShipMap map[string]*ships.Ship

type Sector struct {
	Coords SectorCoords
	Name string

	ShipsByID ShipMap         `spack:"ignore"`
	Listeners *control.ListenerMap `spack:"ignore"`

	manager *SectorManager    `spack:"ignore"`

	chanStop chan int         `spack:"ignore"`
	chanTick <-chan time.Time `spack:"ignore"`
	chanControl chan *ControlCommand `spack:"ignore"`
	chanWarp chan *WarpCommand `spack:"ignore"`

	bodies [1000]simulation.PoweredBody `spack:"ignore"`
}

func (sector *Sector) StorageKey() string {
	return sector.Coords.String()
}

func (sector *Sector) Populate(manager *SectorManager) {
	sector.ShipsByID = make(ShipMap)
	sector.Listeners = control.NewListenerMap()
	sector.manager = manager
	sector.chanStop = make(chan int, 0)
}

func (sector *Sector) Start() {
	sector.LoadShips()
	//sector.DumpShips()
	sector.chanTick = time.Tick(time.Duration(TICK_DELTA) * time.Millisecond)
	sector.chanControl = make(chan *ControlCommand)
	sector.chanWarp = make(chan *WarpCommand)
	go sector.loop()
}

func (sector *Sector) Stop() {
	sector.chanStop <- 1
	<- sector.chanStop
}

func (sector *Sector) Control(session *endpoints.ServerSession, shipID string, join bool, spec control.ControlSpec) (bool, string) {
	var command = &ControlCommand{
		Session: session,
		ShipID: shipID,
		Join: join,
		Spec: spec,
		Reply: make(chan *ControlReply, 1),
	}

	sector.chanControl <- command

	var reply = (<-command.Reply)
	return reply.Success, reply.Error
}

func (sector *Sector) Warp(ship *ships.Ship, in bool) bool {
	var command = &WarpCommand{
		Ship: ship,
		In: in,
		Reply: make(chan *WarpReply, 1),
	}

	sector.chanWarp <- command
	var reply = (<-command.Reply)
	return reply.Success
}

func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v (%d):\n", sector.Name, sector.Coords, len(sector.ShipsByID))
	for _, ship := range sector.ShipsByID {
		ship.Dump()
	}
}
