package sectors

import (
	"drift/common"
	"fmt"
	
)

type Sector struct {
	Coords common.SectorCoords
	Name string
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.Coords.X, sector.Coords.Y)
}

func SectorByCoords(x int, y int) *Sector {
	return &Sector{
		Coords: common.SectorCoords{X: x, Y: y},
	}
}

