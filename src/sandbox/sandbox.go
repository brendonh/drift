package main

import (
	"drift/storage"
	"fmt"
)

type Sector struct {
	X, Y int
	Name string
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.X, sector.Y)
}

func main() {
	client := drift.NewRiakClient("http://localhost:8098")

	sector := Sector{0, 1, "Away"}

	ok := client.Put(&sector)

	if !ok {
		fmt.Printf("Write Failed\n")
		return
	}

	fmt.Printf("Ok\n")

	newSector := Sector{X: 0, Y: 1}
	ok = client.Get(&newSector)

	if !ok {
		fmt.Printf("Read Failed\n")
		return
	}

	fmt.Printf("%s\n", newSector.Name)
}
