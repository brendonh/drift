package sectors

import (
	"drift/ships"
	"drift/simulation"

	"fmt"
)

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

		case command := <-sector.ChanControl:
			sector.control(command)

		case command := <-sector.ChanWarp:
			sector.warp(command)
		}
	}
}


func (sector *Sector) tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location.Body
		pos = pos.EulerIntegrate(1.0)
		ship.Location.Body = pos
	}
}


func (sector *Sector) control(command *ControlCommand) {
	fmt.Printf("%v\n", command)

	var ship, ok = sector.ShipsByID[command.ShipID]
	if !ok {
		command.Reply <- &ControlReply {
			Success: false,
			Error: "No such ship",
		}
		return
	}

	if ship.Owner != command.User.ID() {
		command.Reply <- &ControlReply {
			Success: false,
			Error: "Not your ship",
		}
		return
	}

	command.Reply <- &ControlReply {
		Success: true,
	}

}


func (sector *Sector) warp(command *WarpCommand) {
	// XXX TODO: Out
	var ship = command.Ship

	if ship.Location == nil {
		ship.Location = &ships.ShipLocation {
			ShipID: ship.ID,
			Body: new(simulation.PoweredBody),
		}
	}

	ship.Location.Coords = sector.Coords

	sector.ShipsByID[ship.ID] = ship

	command.Reply <- &WarpReply {
		Success: true,
	}
}

	


	