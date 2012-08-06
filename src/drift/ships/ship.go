package ships

import (
	"drift/storage"
)

type Ship struct {
	ID string
	Owner string `indexed:"true"`
	Name string
}

func (ship *Ship) StorageKey() string {
	return ship.ID
}

func (ship *Ship) SetFromStorageKey(key string) {
	ship.ID = key
}


func NewShip(owner string, name string) *Ship {
	return &Ship{Owner: owner, Name: name}
}


func CreateShip(owner string, name string, client storage.StorageClient) (*Ship, bool) {	
	ship := NewShip(owner, name)
	if !client.Put(ship) {
		return nil, false
	}
	return ship, true
}