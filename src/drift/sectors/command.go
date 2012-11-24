package sectors

import (
	"drift/ships"

	"fmt"
	
	"github.com/brendonh/go-service"
)


type ControlCommand struct {
	User goservice.User
	ShipID string
	Join bool
	Reply chan *ControlReply
}

type ControlReply struct {
	Success bool
	Error string
}

func (control *ControlCommand) String() string {
	return fmt.Sprintf("Control %v by %s of %s", 
		control.Join, control.User.ID(), control.ShipID)
}

// ---------------------------

type WarpCommand struct {
	Ship *ships.Ship
	In bool
	Reply chan *WarpReply
}

type WarpReply struct {
	Success bool
}