package server

import (
	. "drift/common"

	"sync"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/brendonh/go-service"
)

type ServerSession struct {
	id string
	user User
	avatar Entity
	*sync.Mutex
}

// ------------------------------------------
// Session API
// ------------------------------------------

func (session *ServerSession) ID() string {
	return session.id;
}

func (session *ServerSession) User() User {
	return session.user;
}

func (session *ServerSession) SetUser(user User) {
	fmt.Printf("Drift session login: %s (%s)\n", user.DisplayName(), session.id)
	session.user = user
}


// ------------------------------------------
// DriftSession API
// ------------------------------------------

func (session *ServerSession) SetAvatar(entity Entity) {
	fmt.Printf("Session %s (%s) controlling %s\n",
		session.user.DisplayName(), session.id, entity)

	session.avatar = entity
}

func (session *ServerSession) Avatar() Entity {
	return session.avatar
}

func ServerSessionCreator() Session {
	return &ServerSession{
		id: uuid.New(),
		user: nil,
		avatar: nil,
		Mutex: new(sync.Mutex),
	}
}