package sectors

import (
	"drift/ships"
	
	"fmt"
	"time"

	"github.com/brendonh/loge/src/loge"
)

func (sector *Sector) LoadShips() {	
	var db = sector.manager.context.DB()

	var sectorKey = loge.LogeKey(sector.StorageKey())
	var keys = db.Find("shiplocation", "sector", sectorKey)

	fmt.Printf("Loading %d ships...\n", len(keys))
	var start = time.Now()

	var shipsLoaded = 0

	for _, key := range keys {
		var ship = db.ReadOne("ship", loge.LogeKey(key)).(*ships.Ship)
		if ship == nil {
			fmt.Printf("Orphan body: %v\n", key)
			continue
		}
		ship.LoadLocation(db)
		sector.ShipsByID[ship.ID] = ship
		shipsLoaded++
	}

	fmt.Printf("Loaded %d ships in %v\n", shipsLoaded, time.Since(start))
}
