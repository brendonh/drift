package common

import (
	"fmt"
)

type SectorCoords struct {
	X int64
	Y int64
}

func (coords *SectorCoords) String() string {
	return fmt.Sprintf("%d:%d", coords.X, coords.Y)
}