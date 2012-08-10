package server

import (
	. "drift/common"

	"os"
)


type Server struct {
	storage StorageClient
	services API
	endpoints []Endpoint
	
	stopper chan os.Signal
}


func NewServer(
	storage StorageClient, 
	services API) *Server {

	return &Server {
		storage: storage,
		services: services,
		endpoints: make([]Endpoint, 0),
	}
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