package sectors

import (
	. "drift/common"

	"sync"
	"fmt"

	"github.com/brendonh/loge/src/loge"
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

func (manager *SectorManager) Ensure(x int64, y int64) (*Sector, bool) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	var coords = SectorCoords{X: x, Y: y}
	var key = coords.String()

	sector, ok := manager.Sectors[key]
	if !ok {
		fmt.Printf("Loading sector %s\n", key)
		sector = manager.context.DB().ReadOne("sector", loge.LogeKey(key)).(*Sector)

		if sector == nil {
			fmt.Printf("No such sector: %d, %d\n", x, y);
			return nil, false
		}

		sector.Populate(manager)
		manager.Sectors[key] = sector
		sector.Start()
	}
	return sector, true
}


func (manager *SectorManager) Create(x int64, y int64, name string) (*Sector, bool) {
	var coords = SectorCoords{X: x, Y: y}
	var key = coords.String()

	sector, ok := manager.Sectors[key]
	if ok {
		fmt.Printf("Sector exists %s\n", key)
		return sector, false
	}

	var db = manager.context.DB()
	var success = false

	db.Transact(func (t *loge.Transaction) {
		if !t.Exists("sector", loge.LogeKey(key)) {
			sector = &Sector {
				Coords: coords,
				Name: name,
			}
			t.Set("sector", loge.LogeKey(key), sector)
			success = true
		}
	}, 0)

	fmt.Printf("Create success: %v\n", success)

	if success {
		return sector, true
	}

	return nil, false
}


func (manager *SectorManager) Stop() {
	for _, sector := range manager.Sectors {
		sector.Stop()
	}
}