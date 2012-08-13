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

func (p *PoweredBody) RK4Evaluate(dt float64, dIn *derivative, dOut *derivative) {
	var np PoweredBody = PoweredBody{
		Position: p.Position.Add(dIn.Velocity.Muls(dt)),
		Velocity: p.Velocity.Add(dIn.Acceleration.Muls(dt)),
    	Thrust: p.Thrust.Rotate(QMuls(dIn.Spin, dt)),
		Spin: p.Spin,
	}

	dOut.Velocity = np.Velocity
	dOut.Acceleration = np.Acceleration().Muls(dt)
	dOut.Spin = np.Spin
}

func (p *PoweredBody) RK4Integrate(dt float64) *PoweredBody {
	var ds = [5]derivative{}

	p.RK4Evaluate(0.0,      &ds[0], &ds[1])
	p.RK4Evaluate(dt * 0.5, &ds[1], &ds[2])
	p.RK4Evaluate(dt * 0.5, &ds[2], &ds[3])
	p.RK4Evaluate(dt,       &ds[3], &ds[4])

	var dPosition = ds[1].Velocity.Add(
		ds[2].Velocity.Add(ds[3].Velocity).Muls(2.0).Add(ds[4].Velocity).Muls(
		1.0 / 6.0))

	var dVelocity = ds[1].Acceleration.Add(
		ds[2].Acceleration.Add(ds[3].Acceleration).Muls(2.0).Add(ds[4].Acceleration)).Muls(
		1.0 / 6.0)

	var dThrust = QMuls(QAdd(ds[1].Spin, 
		QAdd(QMuls(QAdd(ds[2].Spin, ds[3].Spin), 2.0), ds[4].Spin)), 
		1.0 / 6.0)

	return &PoweredBody{
		Position: p.Position.Add(dPosition),
		Velocity: p.Velocity.Add(dVelocity),
		Thrust: p.Thrust.Rotate(dThrust),
		Spin: p.Spin,
	}
}