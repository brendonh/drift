package common

import (
	"drift/storage"
	"drift/services"
)

type SectorCoords struct {
	X int
	Y int
}

type ServerContext struct {
	StorageClient storage.StorageClient
	Services *services.ServiceCollection
}