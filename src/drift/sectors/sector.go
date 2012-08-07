package sectors

import (
	"fmt"
)

type Sector struct {
	X int
	Y int
	Name string
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.X, sector.Y)
}

func (sector *Sector) SetFromStorageKey(string) {
}

