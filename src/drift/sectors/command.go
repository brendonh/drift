package sectors

import (
	"drift/ships"
	"drift/endpoints"
	"drift/control"

	"fmt"
)


type ControlCommand struct {
	Session *endpoints.ServerSession
	ShipID string
	Join bool
	Spec control.ControlSpec
	Reply chan *ControlReply
}

type ControlReply struct {
	Success bool
	Error string
}

func (control *ControlCommand) String() string {
	return fmt.Sprintf("Control %v by %s of %s", 
		control.Join, control.Session.User().DisplayName(), control.ShipID)
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