package common

import (
	"fmt"
)

type SectorCoords struct {
	X int
	Y int
}

func (coords *SectorCoords) String() string {
	return fmt.Sprintf("%d:%d", coords.X, coords.Y)
}