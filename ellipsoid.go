package globe

import (
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

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
	if radii[0] < 0 || radii[1] < 0 || radii[2] < 0 {
		panic("Ellipsoid with smaller than 0 radii.")
	}

	sqr := radii.Mul(radii)
	return &Ellipsoid{
		radii:               *radii,
		radiiSquared:        *sqr,
		radiiToTheFourth:    *sqr.Mul(sqr),
		oneOverRadiiSquared: vec3.T{1 / sqr[0], 1 / sqr[1], 1 / sqr[2]},
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

// GeodeticSurfaceNormalCart returns the proper surface normal for a point in cartesian coordinates for the ellipsoid.
func (e *Ellipsoid) GeodeticSurfaceNormalCart(posOnEllipsoid *vec3.T) *vec3.T {
	return posOnEllipsoid.Mul(&e.oneOverRadiiSquared).Normalize()
}

// GeodeticSurfaceNormalGeo returns the proper surface normal for a point in geodetic coordinates for the ellipsoid.
func (e *Ellipsoid) GeodeticSurfaceNormalGeo(posOnEllipsoid *Geo3D) *vec3.T {
	lat, lon := posOnEllipsoid.LatLon()
	cosLat := math.Cos(lat)

	return &vec3.T{cosLat * math.Cos(lon), cosLat * math.Sin(lon), math.Sin(lat)}
}
