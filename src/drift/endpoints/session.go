package endpoints

import (
	. "drift/common"

	"sync"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)


type EndpointSession struct {
	id string
	user User
	avatar Entity
	*sync.Mutex
}

func NewEndpointSession() *EndpointSession {
	return &EndpointSession {
		id: uuid.New(),
		user: nil,
		avatar: nil,
		Mutex: new(sync.Mutex),
	}
}


// ------------------------------------------
// Session API
// ------------------------------------------

func (session *EndpointSession) ID() string {
	return session.id;
}

func (session *EndpointSession) User() User {
	return session.user;
}

func (session *EndpointSession) Avatar() Entity {
	return session.avatar;
}

func (session *EndpointSession) SetUser(user User) {
	fmt.Printf("Session login: %s (%s)\n", user.DisplayName(), session.id)
	session.user = user
}

func (session *EndpointSession) SetAvatar(entity Entity) {
	fmt.Printf("Session set entity: %v (%s)\n", entity, session.id)
	session.avatar = entity
}