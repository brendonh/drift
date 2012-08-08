package server

import (
	"drift/common"
	"drift/endpoints"

	"os"
)

type Server struct {
	Context *common.ServerContext
	Endpoints []endpoints.Endpoint
	
	stopper chan os.Signal
}

func NewServer(context *common.ServerContext) *Server {
	httpRpc := endpoints.NewHttpRpcEndpoint(":9999", context)

	return &Server {
		Context: context,
		Endpoints: []endpoints.Endpoint{ httpRpc },
	}
}

func (server *Server) Start() {
	for _, endpoint := range server.Endpoints {
		endpoint.Start()
	}
}

func (server *Server) Stop() {
	for _, endpoint := range server.Endpoints {
		endpoint.Stop()
	}
}

	


