package globe

import "github.com/ungerik/go3d/float64/vec3"

// Ellipsoid holds the necessary information to process a 3d ellipsoidal shape.
type Ellipsoid struct {
	radii               vec3.T
	radiiSquared        vec3.T
	radiiToTheFourth    vec3.T
	oneOverRadiiSquared vec3.T
}

// NewEllipsoid creates a new ellipsoid with the given x, y, and z radii.
func NewEllipsoid(x, y, z float64) *Ellipsoid {
	return NewEllipsoidVec(&vec3.T{x, y, z})
}

// NewEllipsoidVec creates a new ellipsoid with the radii equaling the x, y, and z values of the given vector.
func NewEllipsoidVec(radii *vec3.T) *Ellipsoid {
	if radii.Get(0, 0) < 0 || radii.Get(0, 1) < 0 || radii.Get(0, 2) < 0 {
		panic("Ellipsoid with smaller than 0 radii.")
	}

	sqr := radii.Mul(radii)
	return &Ellipsoid{
		radii:               *radii,
		radiiSquared:        *sqr,
		radiiToTheFourth:    *sqr.Mul(sqr),
		oneOverRadiiSquared: vec3.T{1 / sqr.Get(0, 0), 1 / sqr.Get(0, 1), 1 / sqr.Get(0, 2)},
	}
}

// GetRadii returns the radii vector of the ellipsoid.
func (e *Ellipsoid) GetRadii() *vec3.T {
	return &e.radii
}

// GetRadiiSquared returns the radii vector of the ellipsoid where each radii has been squared.
func (e *Ellipsoid) GetRadiiSquared() *vec3.T {
	return &e.radiiSquared
}

// GetOneOverRadiiSquared returns the radii vector of the ellipsoid where each radii has been squared and then the reciprocal is found.
func (e *Ellipsoid) GetOneOverRadiiSquared() *vec3.T {
	return &e.oneOverRadiiSquared
}
