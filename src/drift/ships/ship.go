package ships

import (
	. "drift/common"

	"fmt"
)

type Ship struct {
	ID string
	Owner string `indexed:"true"`
	Name string
	Location *ShipLocation `msgpack:"-"`
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


func CreateShip(owner string, name string, context DriftServerContext) (*Ship, bool) {	
	var client = context.Storage()
	var id = client.GenerateID()
	ship := NewShip(id, owner, name)
	if !client.Put(ship) {
		return nil, false
	}
	return ship, true
}

func (ship *Ship) SaveLocation(client StorageClient) {
	if ship.Location == nil {
		return;
	}
	
	client.Put(ship.Location)
}

func (ship *Ship) LoadLocation(client StorageClient) bool {
	loc := &ShipLocation{ ShipID: ship.ID }
	if !client.Get(loc) {
		return false
	}
	ship.Location = loc
	return true
}

func (ship *Ship) Dump() {
	fmt.Printf("   %s (%s) (%v, %v, %v, %v)\n", 
		ship.Name,
		ship.ID, 
		ship.Location.Body.Position.String(),
		ship.Location.Body.Velocity.String(),
		ship.Location.Body.Thrust.String(),
		ship.Location.Body.Spin.String())
}