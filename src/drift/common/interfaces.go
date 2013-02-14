package common

import (
	"github.com/brendonh/go-service"
	"github.com/brendonh/loge/src/loge"
)

// ------------------------------------------
// Server
// ------------------------------------------

type DriftServerContext interface {
	goservice.ServerContext
	DB() *loge.LogeDB
}

// ------------------------------------------
// Users and sessions
// ------------------------------------------

type DriftSession interface {
	goservice.Session

	Avatar() Entity
	SetAvatar(Entity)
}


// ------------------------------------------
// Entities
// ------------------------------------------

type Entity interface {
	String() string
}


// ------------------------------------------
// Storage
// ------------------------------------------

type Storable interface {
	StorageKey() string
}

