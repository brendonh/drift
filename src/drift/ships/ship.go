package ships

import (
	"fmt"

	"drift/simulation"

	"github.com/brendonh/loge/src/loge"
)

type Ship struct {
	ID string
	Owner string
	Name string
	Location *simulation.PoweredBody `spack:"ignore"`
}

// ------------------------------------------
// Storage API
// ------------------------------------------

func (ship *Ship) StorageKey() string {
	return ship.ID
}

// ------------------------------------------
// Entity API
// ------------------------------------------

func (ship *Ship) String() string {
	return fmt.Sprintf("%s (%s)", ship.Name, ship.ID)
}

// ------------------------------------------
// Implementation
// ------------------------------------------

func NewShip(id string, owner string, name string) *Ship {	
	return &Ship{ID: id, Owner: owner, Name: name}
}

func (ship *Ship) SaveLocation(db *loge.LogeDB) {
	if ship.Location == nil {
		return;
	}

	db.SetOne("shiplocation", loge.LogeKey(ship.ID), ship.Location)
}

func (ship *Ship) LoadLocation(db *loge.LogeDB) bool {
	var body = db.ReadOne("shiplocation", loge.LogeKey(ship.ID)).(*simulation.PoweredBody)
	if body == nil {
		return false
	}
	ship.Location = body
	return true
}

func (ship *Ship) Dump() {
	fmt.Printf("   %s (%s) (%v, %v, %v, %v)\n", 
		ship.Name,
		ship.ID, 
		ship.Location.Position.String(),
		ship.Location.Velocity.String(),
		ship.Location.Thrust.String(),
		ship.Location.Spin.String())
}


func GetShip(body *simulation.PoweredBody, db *loge.LogeDB) *Ship {
	var ship = db.ReadOne("ship", loge.LogeKey(body.ShipID)).(*Ship)
	if ship == nil {
		return nil
	}
	ship.Location = body
	return ship
}