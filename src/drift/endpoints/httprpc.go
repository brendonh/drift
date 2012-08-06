package endpoints

import (
	"drift/services"

	"fmt"
	"net"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
)

type HttpRpcEndpoint struct {
	Address string
	listener net.Listener
	
	// XXX BGH TODO: General context here
	ServiceCollection *services.ServiceCollection
}


func NewHttpRpcEndpoint(address string, sc *services.ServiceCollection) HttpRpcEndpoint {
	return HttpRpcEndpoint{ 
		Address: address,
		ServiceCollection: sc,
	}
}

func (endpoint *HttpRpcEndpoint) Start() bool {
	if endpoint.listener != nil {
		return false
	}

	listener, error := net.Listen("tcp", endpoint.Address)
	if error != nil {
		fmt.Printf("Error starting HTTP RPC endpoint: %v\n", error)
		return false
	}

	endpoint.listener = listener
	
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.Handle("/", endpoint)
	go http.Serve(listener, mux)

	return true
}

func (endpoint *HttpRpcEndpoint) Stop() bool {
	if endpoint.listener == nil {
		return true
	}

	if error := endpoint.listener.Close(); error != nil {
		fmt.Printf("Error stopping HTTP RPC endpoint: %v\n", error)
		return false
	}

	endpoint.listener = nil
	return true
}

func (endpoint *HttpRpcEndpoint) ServeHTTP(response http.ResponseWriter, req *http.Request) {

	bits := strings.SplitN(req.URL.Path[1:], "/", 2)

	if len(bits) != 2 {
		http.NotFound(response, req)
		return
	}

	req.ParseForm()

	var form = make(services.APIData)
	for k, v := range req.Form {
		form[k] = v[0]
	}

	ok, errors, resp := endpoint.ServiceCollection.HandleCall(bits[0], bits[1], form)

	if errors != nil {
		response.WriteHeader(400)
	}

	response.Header().Add("Content-Type", "text/plain")

	jsonReply, _ := json.Marshal(services.Response(ok, errors, resp))
	response.Header().Add("Content-Length", strconv.Itoa(len(jsonReply)))

	response.Write(jsonReply)
	

}