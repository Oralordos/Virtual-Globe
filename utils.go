package globe

import (
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

func rotateAroundAxis(vec, axis *vec3.T, theta float64) *vec3.T {
	x, y, z := vec[0], vec[1], vec[2]
	u, v, w := axis[0], axis[1], axis[2]

	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)

	ms := axis.LengthSqr()
	m := math.Sqrt(ms)

	return &vec3.T{
		((u * (u*x + v*y + w*z)) + (((x * (v*v + w*w)) - (u * (v*y + w*z))) * cosTheta) + (m * ((-w * y) + (v * z)) * sinTheta)) / ms,
		((v * (u*x + v*y + w*z)) + (((y * (u*u + w*w)) - (v * (u*x + w*z))) * cosTheta) + (m * ((w * x) - (u * z)) * sinTheta)) / ms,
		((w * (u*x + v*y + w*z)) + (((z * (u*u + v*v)) - (w * (u*x + v*y))) * cosTheta) + (m * (-(v * x) + (u * y)) * sinTheta)) / ms,
	}
}
