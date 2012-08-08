package simulation

import (
	"drift/common"
	. "github.com/brendonh/s3dm-go"
)


var MAX_SPEED_SQUARED = common.MAX_SPEED * common.MAX_SPEED

type PoweredBody struct {
	Position V3
	Velocity V3
	Thrust V3
}

func (p *PoweredBody) Acceleration() *V3 {
	var scale = p.Thrust.Dot(&p.Velocity) / MAX_SPEED_SQUARED
	return p.Thrust.Sub(p.Velocity.Muls(scale))
}

// ------------------------------------------
// RK4
// ------------------------------------------

type derivative struct {
	Velocity V3
	Acceleration V3
}

func (p *PoweredBody) RK4Evaluate(dt float64, derivativeIn *derivative) *derivative {
	var np *PoweredBody = &PoweredBody{
		Position: *p.Position.Add(derivativeIn.Velocity.Muls(dt)),
		Velocity: *p.Velocity.Add(derivativeIn.Acceleration.Muls(dt)),
    	Thrust: p.Thrust,
	}
	var derivativeOut = &derivative {
		Velocity: np.Velocity,
		Acceleration: *np.Acceleration().Muls(dt),
	}
	return derivativeOut
}

func (p *PoweredBody) RK4Integrate(dt float64) *PoweredBody {
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

	return &PoweredBody{
		Position: *p.Position.Add(dPosition),
		Velocity: *p.Velocity.Add(dVelocity),
		Thrust: p.Thrust,
	}
}