package sectors

import (
	. "drift/common"

	"sync"
	"fmt"
)

// ------------------------------------------
// Manager
// ------------------------------------------

type SectorManager struct {
	Sectors map[string]*Sector
	context DriftServerContext
	*sync.Mutex
}

func NewSectorManager(context DriftServerContext) *SectorManager {
	return &SectorManager {
		Sectors: make(map[string]*Sector),
		context: context,
		Mutex: new(sync.Mutex),
	}
}

func (manager *SectorManager) Ensure(x int, y int) (*Sector, bool) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	
	var coords = SectorCoords{X: x, Y: y}
	var key = coords.String()

	sector, ok := manager.Sectors[key]
	if !ok {
		fmt.Printf("Loading sector %s\n", key)
		sector = SectorByCoords(x, y, manager)
		if !manager.context.Storage().Get(sector) {
			fmt.Printf("No such sector: %d, %d\n", x, y);
			return nil, false
		}
		manager.Sectors[key] = sector
		sector.Start()
	}
	return sector, true
}


func (manager *SectorManager) Create(x int, y int, name string) (*Sector, bool) {
	var coords = SectorCoords{X: x, Y: y}
	var key = coords.String()

	sector, ok := manager.Sectors[key]
	if ok {
		fmt.Printf("Sector exists %s\n", key)
		return sector, false
	}

	sector = SectorByCoords(x, y, manager)
	if manager.context.Storage().Get(sector) {
		fmt.Printf("Sector exists: %d, %d\n", x, y);
		return sector, false
	}

	sector.Name = name
	manager.context.Storage().Put(sector)
	return sector, true
	
}


func (manager *SectorManager) Stop() {
	for _, sector := range manager.Sectors {
		sector.Stop()
	}
}