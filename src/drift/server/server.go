package server

import (
	"drift/sectors"
	"drift/endpoints"

	"os"

	. "github.com/brendonh/go-service"
	"github.com/brendonh/loge/src/loge"
)


type DriftServer struct {
	Server

	//storage StorageClient
	db *loge.LogeDB
	SectorManager *sectors.SectorManager

	stopper chan os.Signal
}
	

func NewDriftServer(
	db *loge.LogeDB, 
	services API) *DriftServer {

	var server = &DriftServer {
		*NewServer(services, endpoints.ServerSessionCreator),
		db,
		nil,
		nil,
	}

	server.SectorManager = sectors.NewSectorManager(server)
	return server

}

// ------------------------------------------
// Context API
// ------------------------------------------

func (server *DriftServer) DB() *loge.LogeDB {
	return server.db
}