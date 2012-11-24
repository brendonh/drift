package sectors

import (
	. "drift/common"
	"drift/ships"
	
	"fmt"
	"time"
)

func (sector *Sector) loadShips() {	
	var client = sector.manager.context.Storage()

	searchLoc := &ships.ShipLocation{ Coords: sector.Coords }
	foundLocs := make([]ships.ShipLocation, 0)
	client.IndexLookup(searchLoc, &foundLocs, "Coords")

	var chunkSize = 32
	var totalShips = len(foundLocs)
	fmt.Printf("Loading %d ships...\n", totalShips)
	var start = time.Now()

	var chunkAcks = make(chan int)
	var shipStream = make(chan *ships.Ship)

	var chunks = 0

	for chunk := 0; chunk <= totalShips / chunkSize; chunk++ {
		var start = chunk * chunkSize
		var end = (chunk + 1) * chunkSize

		if end > totalShips {
			end = totalShips
		}

		if end <= start {
			break
		}

		go sector.loadShipChunk(foundLocs[start:end], client, chunkAcks, shipStream)
		chunks++;
	}

	var finishedChunks = 0
	var shipsLoaded = 0
	var ship *ships.Ship

	for finishedChunks < chunks {
		select {
		case <-chunkAcks:
			finishedChunks += 1
		case ship = <-shipStream:
			shipsLoaded += 1
			sector.ShipsByID[ship.ID] = ship
		}
	}

	fmt.Printf("Loaded %d ships in %v\n", shipsLoaded, time.Since(start))
}

func (sector *Sector) loadShipChunk(
	locs []ships.ShipLocation, client StorageClient, 
	done chan int, stream chan *ships.Ship) {
	for i := range locs {
		var loc = locs[i]
		var ship = loc.GetShip(client)
		if ship == nil {
			fmt.Printf("Orphan location: %v\n", loc)
			client.Delete("ShipLocation", loc.StorageKey())
			continue
		}
		stream <- ship
	}
	done <- 1
}
