package ships

import (
	"drift/storage"
)

type Ship struct {
	ID string
	Owner string `indexed:"true"`
	Name string
	Location *ShipLocation `msgpack:"-"`
}

func (ship *Ship) StorageKey() string {
	return ship.ID
}

func NewShip(id string, owner string, name string) *Ship {	
	return &Ship{ID: id, Owner: owner, Name: name}
}


func CreateShip(owner string, name string, client storage.StorageClient) (*Ship, bool) {	
	var id = client.GenerateID()
	ship := NewShip(id, owner, name)
	if !client.Put(ship) {
		return nil, false
	}
	return ship, true
}

func (ship *Ship) SaveLocation(client storage.StorageClient) {
	if ship.Location == nil {
		return;
	}
	
	client.Put(ship.Location)
}

func (ship *Ship) LoadLocation(client storage.StorageClient) bool {
	loc := &ShipLocation{ ShipID: ship.ID }
	if !client.Get(loc) {
		return false
	}
	ship.Location = loc
	return true
}