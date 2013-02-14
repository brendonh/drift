package sectors

import (
	"drift/simulation"

	"fmt"
	"time"
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
			var start = time.Now()
			sector.tick()
			fmt.Printf("Tick: %v\n", time.Since(start))

		case command := <-sector.chanControl:
			sector.control(command)

		case command := <-sector.chanWarp:
			sector.warp(command)
		}
	}
}


func (sector *Sector) tick() {
	for _, ship := range sector.ShipsByID {
		var pos = ship.Location
		pos = pos.EulerIntegrate(1.0)
		ship.Location = pos
	}
}


func (sector *Sector) control(command *ControlCommand) {
	var ship, ok = sector.ShipsByID[command.ShipID]
	if !ok {
		command.Reply <- &ControlReply {
			Success: false,
			Error: "No such ship",
		}
		return
	}

	var user = command.Session.User()
	if ship.Owner != user.ID() {
		command.Reply <- &ControlReply {
			Success: false,
			Error: "Not your ship",
		}
		return
	}

	sector.Listeners.Set(ship, command.Session, command.Spec)

	fmt.Printf("%s controlling %s (%v)\n", user.DisplayName(), ship.Name, command.Spec)

	command.Reply <- &ControlReply {
		Success: true,
	}

}


func (sector *Sector) warp(command *WarpCommand) {
	// XXX TODO: Out
	var ship = command.Ship

	if ship.Location == nil {
		ship.Location = &simulation.PoweredBody{
			ShipID: ship.ID,
			Coords: sector.Coords,
		}
	}

	sector.ShipsByID[ship.ID] = ship

	command.Reply <- &WarpReply {
		Success: true,
	}
}

	


	