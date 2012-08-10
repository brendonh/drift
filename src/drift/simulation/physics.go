package simulation

import (
	"drift/common"
	. "github.com/klkblake/s3dm"
)


var MAX_SPEED_SQUARED = common.MAX_SPEED * common.MAX_SPEED

type PoweredBody struct {
	Position V3
	Velocity V3
	Thrust V3
	Spin Qtrnn
}

func (p *PoweredBody) Acceleration() V3 {
	var scale = p.Thrust.Dot(p.Velocity) / MAX_SPEED_SQUARED
	return p.Thrust.Sub(p.Velocity.Muls(scale))
}


// ------------------------------------------
// RK4
// ------------------------------------------

type derivative struct {
	Velocity V3
	Acceleration V3
	Spin Qtrnn
}

func QMuls (q Qtrnn, s float64) Qtrnn {
	return Qtrnn{q.X * s, q.Y * s, q.Z * s, q.W * s}
}

func QAdd(q1 Qtrnn, q2 Qtrnn) Qtrnn {
	return Qtrnn{q1.X + q2.X, q1.Y + q2.Y, q1.Z + q2.Z, q1.W + q2.W}
}

func (p *PoweredBody) RK4Evaluate(dt float64, derivativeIn *derivative) *derivative {
	var np *PoweredBody = &PoweredBody{
		Position: p.Position.Add(derivativeIn.Velocity.Muls(dt)),
		Velocity: p.Velocity.Add(derivativeIn.Acceleration.Muls(dt)),
    	Thrust: p.Thrust.Rotate(QMuls(derivativeIn.Spin, dt)),
		Spin: p.Spin,
	}

	var derivativeOut = &derivative {
		Velocity: np.Velocity,
		Acceleration: np.Acceleration().Muls(dt),
		Spin: np.Spin,
	}
	return derivativeOut
}

func (p *PoweredBody) RK4Integrate(dt float64) *PoweredBody {
	var a = p.RK4Evaluate(0.0, &derivative{})
	var b = p.RK4Evaluate(dt * 0.5, a)
	var c = p.RK4Evaluate(dt * 0.5, b)
	var d = p.RK4Evaluate(dt, c)

	var dPosition = a.Velocity.Add(
		b.Velocity.Add(c.Velocity).Muls(2.0).Add(d.Velocity).Muls(
		1.0 / 6.0))

	var dVelocity = a.Acceleration.Add(
		b.Acceleration.Add(c.Acceleration).Muls(2.0).Add(d.Acceleration)).Muls(
		1.0 / 6.0)

	var dThrust = QMuls(QAdd(a.Spin, 
		QAdd(QMuls(QAdd(b.Spin, c.Spin), 2.0), d.Spin)), 
		1.0 / 6.0)

	return &PoweredBody{
		Position: p.Position.Add(dPosition),
		Velocity: p.Velocity.Add(dVelocity),
		Thrust: p.Thrust.Rotate(dThrust),
		Spin: p.Spin,
	}
}