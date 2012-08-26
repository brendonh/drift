package endpoints

import (
	. "drift/common"

	"fmt"

	. "github.com/brendonh/go-service"
)


type EndpointSession struct {
	BasicSession
	avatar Entity
}

func NewEndpointSession() *EndpointSession {
	return &EndpointSession {
		*NewBasicSession(),
		nil,
	}
}


// ------------------------------------------
// Session API
// ------------------------------------------

func (session *EndpointSession) Avatar() Entity {
	return session.avatar;
}

func (session *EndpointSession) SetAvatar(entity Entity) {
	fmt.Printf("Session set entity: %v (%s)\n", entity, session.ID())
	session.avatar = entity
}