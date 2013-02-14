package simulation

import (

)

type SectorState struct {
	Bodies []PoweredBody
}

func NewSectorState() *SectorState {
	return &SectorState {
		Bodies: make([]PoweredBody, 0, 100),
	}
}

func (state *SectorState) AddBody(body PoweredBody) *SectorState {
	return &SectorState{
		Bodies: append(state.Bodies, body),
	}
}

func (state *SectorState) RemoveBody(shipID string) *SectorState {
	var index = -1
	var i int
	for i, b := range state.Bodies {
		if b.ShipID == shipID {
			index = i
			break
		}
	}

	var newBodies []PoweredBody
	if index >= 0 {
		newBodies = make([]PoweredBody, 0, len(state.Bodies))
		copy(newBodies, state.Bodies[:i])
		copy(newBodies, state.Bodies[i+1:])
	}

	return &SectorState{
		Bodies: newBodies,
	}
}