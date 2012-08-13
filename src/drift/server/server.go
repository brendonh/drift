package server

import (
	. "drift/common"
	"drift/sectors"

	"os"
)


type Server struct {
	storage StorageClient
	services API
	endpoints []Endpoint

	SectorManager *sectors.SectorManager
	
	stopper chan os.Signal
}


func NewServer(
	storage StorageClient, 
	services API) *Server {

	var server = &Server {
		storage: storage,
		services: services,
		endpoints: make([]Endpoint, 0),
	}

	server.SectorManager = sectors.NewSectorManager(server)
	return server

}

func (server *Server) AddEndpoint(endpoint Endpoint) {
	server.endpoints = append(server.endpoints, endpoint)
}

func (server *Server) Start() {
	for _, endpoint := range server.endpoints {
		endpoint.Start()
	}
}

func (server *Server) Stop() {
	for _, endpoint := range server.endpoints {
		endpoint.Stop()
	}
	server.SectorManager.Stop()
}


// ------------------------------------------
// Context API
// ------------------------------------------

func (server *Server) Storage() StorageClient {
	return server.storage
}

func (server *Server) API() API {
	return server.services
}