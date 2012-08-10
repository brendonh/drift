package endpoints

import (
	. "drift/common"
	//"drift/services"

	"fmt"
	"reflect"
	"bytes"
	"net"
	"net/http"

	"github.com/ugorji/go-msgpack"
	"code.google.com/p/go.net/websocket"
)

type WebsocketEndpoint struct {
	Address string
	listener net.Listener
	context ServerContext
}


func NewWebsocketEndpoint(address string, context ServerContext) Endpoint {
	return &WebsocketEndpoint{
		Address: address,
		context: context,
	}
}

func (endpoint *WebsocketEndpoint) Start() bool {
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

	var handler = func(ws *websocket.Conn) {
		endpoint.Handle(ws)
	}

	mux.Handle("/", websocket.Handler(handler))
	go http.Serve(listener, mux)

	return true
}


func (endpoint *WebsocketEndpoint) Stop() bool {
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


const (
	APIFrame = 'a'
	PositionFrame = 'p'
	PingFrame = 'P'
)

func (endpoint *WebsocketEndpoint) Handle(ws *websocket.Conn) {
	ws.PayloadType = websocket.BinaryFrame

	var buf = make([]byte, 1024 * 64)

	for {
		msgLength, err := ws.Read(buf)
		
		if err != nil {
			fmt.Printf("WS error: %v\n", err)
			break
		}

		if msgLength == 0 {
			continue
		}

		switch buf[0] {
		case APIFrame:
			go endpoint.HandleAPI(buf[1:msgLength], ws)
		case PositionFrame:
			fmt.Printf("Position frame: %v\n", buf[1:msgLength])
		default:
			fmt.Printf("Unknown frame: %v\n", buf[:msgLength])
		}
	}
}


func (endpoint *WebsocketEndpoint) HandleAPI(buf []byte, ws *websocket.Conn) {
	var data APIData
	var resolver = msgpack.DefaultDecoderContainerResolver
	resolver.MapType = reflect.TypeOf(make(APIData))

	var dec = msgpack.NewDecoder(bytes.NewReader(buf), &resolver)
	
	var err = dec.Decode(&data)

	if err != nil {
		fmt.Printf("Decode err: %v\n", err)
		return
	}

	var response = endpoint.context.API().HandleRequest(data, endpoint.context)

	if id, ok := data["id"]; ok {
		response["id"] = id
	}

	reply, err := msgpack.Marshal(response)

	if err != nil {
		fmt.Printf("Encode err: %#v\n", err)
		return
	}

	ws.Write(reply)
}