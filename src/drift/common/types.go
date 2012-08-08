package common

import (
	"drift/storage"
)

type SectorCoords struct {
	X int
	Y int
}

type ServerContext struct {
	StorageClient storage.StorageClient
	Services API
}


// ------------------------------------------
// API
// ------------------------------------------


type APIArg struct {
	Name string
	ArgType int
	Required bool
	Default interface{}
	Extra interface{}
}

type APIMethod struct {
	Name string
	ArgSpec []APIArg
	Handler APIHandler
}

type APIData map[string]interface{}

type APIHandler func(APIData, *ServerContext) (bool, APIData)


// ------------------------------------------
// Services
// ------------------------------------------

type APIService interface {
	Name() string
	AddMethod(string, []APIArg, APIHandler)
	FindMethod(string) *APIMethod
}

type API interface {
	AddService(APIService)
	HandleRequest(APIData, *ServerContext) APIData
	HandleCall(string, string, APIData, *ServerContext) (bool, []string, APIData)
}