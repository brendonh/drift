package simulation

import (
	. "github.com/brendonh/s3dm-go"
)

var MAX_SPEED = 50.0
var MAX_SPEED_SQUARED = MAX_SPEED * MAX_SPEED

type Powered struct {
	Position V3
	Velocity V3
	Orientation V3
	ThrustAccel float64
}

func (p *Powered) Acceleration() *V3 {
	var scale = p.Orientation.Dot(&p.Velocity) / MAX_SPEED_SQUARED
	var acceleration = p.Orientation.Sub(p.Velocity.Muls(scale))
	return acceleration.Muls(p.ThrustAccel)
}

// ------------------------------------------
// RK4
// ------------------------------------------

type derivative struct {
	Velocity V3
	Acceleration V3
}

func (p *Powered) RK4Evaluate(dt float64, derivativeIn *derivative) *derivative {
	var np *Powered = &Powered{
		Position: *p.Position.Add(derivativeIn.Velocity.Muls(dt)),
		Velocity: *p.Velocity.Add(derivativeIn.Acceleration.Muls(dt)),
		Orientation: p.Orientation,
		ThrustAccel: p.ThrustAccel,
	}
	var derivativeOut = &derivative {
		Velocity: np.Velocity,
		Acceleration: *np.Acceleration().Muls(dt),
	}
	return derivativeOut
}

func (p *Powered) RK4Integrate(dt float64) *Powered {
	var a = p.RK4Evaluate(0.0, &derivative{})
	var b = p.RK4Evaluate(dt * 0.5, a)
	var c = p.RK4Evaluate(dt * 0.5, b)
	var d = p.RK4Evaluate(dt, c)

	var dPosition = a.Velocity.Add(
		b.Velocity.Add(&c.Velocity).Muls(2.0).Add(&d.Velocity).Muls(
		1.0 / 6.0))

	var dVelocity = a.Acceleration.Add(
		b.Acceleration.Add(&c.Acceleration).Muls(2.0).Add(&d.Acceleration)).Muls(
		1.0 / 6.0)

	return &Powered{
		Position: *p.Position.Add(dPosition),
		Velocity: *p.Velocity.Add(dVelocity),
		Orientation: p.Orientation,
		ThrustAccel: p.ThrustAccel,
	}
}