package server

import (
	. "drift/common"
	"drift/sectors"

	"os"

	. "github.com/brendonh/go-service"
)


type DriftServer struct {
	Server

	storage StorageClient
	SectorManager *sectors.SectorManager

	stopper chan os.Signal
}


func NewDriftServer(
	storage StorageClient, 
	services API) *DriftServer {

	var server = &DriftServer {
		*NewServer(services),
		storage,
		nil,
		nil,
	}

	server.SectorManager = sectors.NewSectorManager(server)
	return server

}


// ------------------------------------------
// Context API
// ------------------------------------------

func (server *DriftServer) Storage() StorageClient {
	return server.storage
}
