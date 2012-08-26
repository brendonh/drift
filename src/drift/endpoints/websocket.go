package endpoints

import (
	"fmt"

	. "github.com/brendonh/go-service"

	"code.google.com/p/go.net/websocket"
)


const (
	APIFrame = 'a'
	PositionFrame = 'p'
	PingFrame = 'P'
)

func DriftMessageHandler(endpoint *WebsocketEndpoint, buf []byte, session Session, conn *websocket.Conn) {
	switch buf[0] {
	case APIFrame:
		go endpoint.HandleAPI(buf[1:], session, conn)
	case PositionFrame:
		fmt.Printf("Position frame: %v\n", buf[1:])
	default:
		fmt.Printf("Unknown frame: %v\n", buf)
	}
}
