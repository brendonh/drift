package sectors

import (
	. "drift/common"
	"drift/ships"
	"drift/simulation"
	
	"fmt"
	"time"

	"github.com/brendonh/go-service"
)


type ShipMap map[string]*ships.Ship

type Sector struct {
	Coords SectorCoords
	Name string

	ShipsByID ShipMap         `msgpack:"-"`
	ChanControl chan *ControlCommand `msgpack:"-"`
	ChanWarp chan *WarpCommand `msgpack:"-"`

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
	sector.ChanControl = make(chan *ControlCommand)
	sector.ChanWarp = make(chan *WarpCommand)
	go sector.loop()
}

func (sector *Sector) Stop() {
	sector.chanStop <- 1
	<- sector.chanStop
}

func (sector *Sector) Control(user goservice.User, shipID string, join bool) (bool, string) {
	
	var command = &ControlCommand{
		User: user,
		ShipID: shipID,
		Join: join,
		Reply: make(chan *ControlReply, 1),
	}

	sector.ChanControl <- command

	var reply = (<-command.Reply)
	return reply.Success, reply.Error
}

func (sector *Sector) Warp(ship *ships.Ship, in bool) bool {
	var command = &WarpCommand{
		Ship: ship,
		In: in,
		Reply: make(chan *WarpReply, 1),
	}

	sector.ChanWarp <- command
	var reply = (<-command.Reply)
	return reply.Success
}

func (sector *Sector) DumpShips() {
	fmt.Printf("Ships in %s %v (%d):\n", sector.Name, sector.Coords, len(sector.ShipsByID))
	for _, ship := range sector.ShipsByID {
		ship.Dump()
	}
}
