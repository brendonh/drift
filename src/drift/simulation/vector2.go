package simulation

import (
	"math"
	"fmt"
)


type V2 struct {
	X, Y float64
}

func (v V2) Equals(o V2) bool {
	return v.X == o.X && v.Y == o.Y
}

func (v V2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v V2) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v V2) Dot(o V2) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v V2) Reflect(norm V2) V2 {
	distance := 2.0 * v.Dot(norm)
	return V2{
		v.X - distance*norm.X,
		v.Y - distance*norm.Y,
	}
}

func (v V2) Perp() V2 {
	return V2 { v.Y, -v.X }
}

func (v V2) SetLength(l float64) V2 {
	return v.Unit().Muls(l)
}

func (v V2) Unit() V2 {
	return v.Muls(1 / v.Length())
}

func (v V2) Add(o V2) V2 {
	return V2{
		v.X + o.X,
		v.Y + o.Y,
	}
}

func (v V2) Adds(o float64) V2 {
	return v.Add(V2{o, o})
}

func (v V2) Sub(o V2) V2 {
	return V2{
		v.X - o.X,
		v.Y - o.Y,
	}
}

func (v V2) Subs(o float64) V2 {
	return v.Sub(V2{o, o})
}

func (v V2) Mul(o V2) V2 {
	return V2{
		v.X * o.X,
		v.Y * o.Y,
	}
}

func (v V2) Muls(o float64) V2 {
	return v.Mul(V2{o, o})
}


func (v V2) Rotate(rot V2) V2 {
	return V2 {
		(v.X * rot.X) - (v.Y * rot.Y),
		(v.X * rot.Y) + (v.Y * rot.X),
	}
}

func (v V2) String() string {
	return fmt.Sprintf("<%.3f, %.3f>", v.X, v.Y)
}
