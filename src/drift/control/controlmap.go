package control

import (
	"drift/ships"
	"drift/endpoints"

	"container/list"
)

type LMShipMap map[*ships.Ship]*list.List
type LMSessionMap map[*endpoints.ServerSession]*list.List

type ControlSpec bool

type ListenerMap struct {
	byShip LMShipMap
	bySession LMSessionMap
}

func NewListenerMap() *ListenerMap {
	return &ListenerMap {
		make(LMShipMap, 64),
		make(LMSessionMap, 64),
	}
}

type Control struct {
	Ship *ships.Ship
	Session *endpoints.ServerSession
	Spec ControlSpec
}

func NewControl(ship *ships.Ship, session *endpoints.ServerSession, spec ControlSpec) *Control {
	return &Control{
		Ship: ship,
		Session: session,
		Spec: spec,
	}
}


func (lm *ListenerMap) Set(ship *ships.Ship, session *endpoints.ServerSession, spec ControlSpec) {
	lm.Remove(ship, session)

	var control = NewControl(ship, session, spec)

	shipMap, ok := lm.byShip[ship]
	if !ok {
		shipMap = list.New()
		lm.byShip[ship] = shipMap
	}
	shipMap.PushBack(control)

	sessionMap, ok := lm.bySession[session]
	if !ok {
		sessionMap = list.New()
		lm.bySession[session] = sessionMap
	}
	sessionMap.PushBack(control)
}

func (lm *ListenerMap) Remove(ship *ships.Ship, session *endpoints.ServerSession) {
	lm.byShip.RemoveSession(ship, session)
	lm.bySession.RemoveShip(ship, session)
}

func (lm *ListenerMap) ClearShip(ship *ships.Ship) {
	var shipMap, ok = lm.byShip[ship]
	if !ok { return }
	for e := shipMap.Front(); e != nil; e = e.Next() {
		var session = e.Value.(*Control).Session
		lm.bySession.RemoveShip(ship, session)
	}
	delete(lm.byShip, ship)
}

func (lm *ListenerMap) ClearSession(session *endpoints.ServerSession) {
	var sessionMap, ok = lm.bySession[session]
	if !ok { return }
	for e := sessionMap.Front(); e != nil; e = e.Next() {
		var ship = e.Value.(*Control).Ship
		lm.byShip.RemoveSession(ship, session)
	}
	delete(lm.bySession, session)
}


// -------------------------------
// Queries
// -------------------------------

type SessionCallback func(*endpoints.ServerSession, ControlSpec)
type ShipCallback func(*ships.Ship, ControlSpec)

func (lm *ListenerMap) SessionApply(ship *ships.Ship, callback SessionCallback) {
	var shipMap, ok = lm.byShip[ship]
	if !ok { return }

	for e := shipMap.Front(); e != nil; e = e.Next() {
		var control = e.Value.(*Control)
		callback(control.Session, control.Spec)
	}
}


func (lm *ListenerMap) ShipApply(session *endpoints.ServerSession, callback ShipCallback) {
	var sessionMap, ok = lm.bySession[session]
	if !ok { return }

	for e := sessionMap.Front(); e != nil; e = e.Next() {
		var control = e.Value.(*Control)
		callback(control.Ship, control.Spec)
	}
}

func (lm *ListenerMap) HasControl(ship *ships.Ship, session *endpoints.ServerSession) bool {
	var sessionMap, ok = lm.bySession[session]
	if !ok { return false }

	for e := sessionMap.Front(); e != nil; e = e.Next() {
		var control = e.Value.(*Control)
		if (control.Ship == ship) {
			return bool(control.Spec)
		}
	}
	return false
}


// -------------------------------
// Internal
// -------------------------------

func (shipMap LMShipMap) RemoveSession(ship *ships.Ship, session *endpoints.ServerSession) {
	var shipList, ok = shipMap[ship]
	if !ok {
		return;
	}

	for e := shipList.Front(); e != nil; e = e.Next() {
		if e.Value.(*Control).Session == session {
			shipList.Remove(e)
		}
	}

	if shipList.Front() == nil {
		delete(shipMap, ship)
	}
}


func (sessionMap LMSessionMap) RemoveShip(ship *ships.Ship, session *endpoints.ServerSession) {
	var sessionList, ok = sessionMap[session]
	if !ok {
		return;
	}

	for e := sessionList.Front(); e != nil; e = e.Next() {
		if e.Value.(*Control).Ship == ship {
			sessionList.Remove(e)
		}
	}

	if sessionList.Front() == nil {
		delete(sessionMap, session)
	}
}
