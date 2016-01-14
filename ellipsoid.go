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

	sqr := vec3.Mul(radii, radii)
	return &Ellipsoid{
		radii:               *radii,
		radiiSquared:        sqr,
		radiiToTheFourth:    vec3.Mul(&sqr, &sqr),
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
	v := vec3.Mul(posOnEllipsoid, &e.oneOverRadiiSquared)
	return v.Normalize()
}

// GeodeticSurfaceNormalGeo returns the proper surface normal for a point in geodetic coordinates for the ellipsoid.
func (e *Ellipsoid) GeodeticSurfaceNormalGeo(posOnEllipsoid *Geo3D) *vec3.T {
	p := posOnEllipsoid.ToRad()
	lat, lon := p.LatLon()
	cosLat := math.Cos(lat)

	return &vec3.T{cosLat * math.Cos(lon), cosLat * math.Sin(lon), math.Sin(lat)}
}

// ToVec3From2D converts 2D geodetic coordinates to cartesian coordinates.
func (e *Ellipsoid) ToVec3From2D(geo *Geo2D) *vec3.T {
	return e.ToVec3(&Geo3D{*geo, 0})
}

// ToVec3 converts 3D geodetic coordinates to cartesian coordinates.
func (e *Ellipsoid) ToVec3(geo *Geo3D) *vec3.T {
	n := e.GeodeticSurfaceNormalGeo(geo)
	k := vec3.Mul(&e.radiiSquared, n)
	gamma := math.Sqrt(k[0]*n[0] + k[1]*n[1] + k[2]*n[2])

	surface := k.Scale(1 / gamma)
	return surface.Add(n.Scale(geo.GetHeight()))
}

// ToGeo2D takes a point on the surface of an ellipsoid, and returns the latitude and longitude.
func (e *Ellipsoid) ToGeo2D(posOnEllipsoid *vec3.T) *Geo2D {
	n := e.GeodeticSurfaceNormalCart(posOnEllipsoid)
	return &Geo2D{math.Asin(n[2] / n.Length()), math.Atan2(n[1], n[0]), true}
}

// ScaleToGeocentricSurface converts an arbitrary 3D point to a point on the ellipsoid, going directly towards the center of the ellipsoid.
func (e *Ellipsoid) ScaleToGeocentricSurface(pos *vec3.T) *vec3.T {
	beta := 1 / math.Sqrt(
		(pos[0]*pos[0])*e.oneOverRadiiSquared[0]+
			(pos[1]*pos[1])*e.oneOverRadiiSquared[1]+
			(pos[2]*pos[2])*e.oneOverRadiiSquared[2],
	)
	v := pos.Scaled(beta)
	return &v
}

// ScaleToGeodeticSurface converts an arbitrary 3D point to the closest point on the ellipsoid.
func (e *Ellipsoid) ScaleToGeodeticSurface(pos *vec3.T) *vec3.T {
	beta := 1 / math.Sqrt(
		(pos[0]*pos[0])*e.oneOverRadiiSquared[0]+
			(pos[1]*pos[1])*e.oneOverRadiiSquared[1]+
			(pos[2]*pos[2])*e.oneOverRadiiSquared[2],
	)
	nP := &vec3.T{
		beta * pos[0] * e.oneOverRadiiSquared[0],
		beta * pos[1] * e.oneOverRadiiSquared[1],
		beta * pos[2] * e.oneOverRadiiSquared[2],
	}
	n := nP.Length()
	alpha := (1 - beta) * (pos.Length() / n)

	x2 := pos[0] * pos[0]
	y2 := pos[1] * pos[1]
	z2 := pos[2] * pos[2]

	da := 0.0
	db := 0.0
	dc := 0.0

	s := 0.0
	dSdA := 1.0
	first := true

	for math.Abs(s) > 1e-10 || first {
		first = false
		alpha -= (s / dSdA)

		da = 1 + (alpha * e.oneOverRadiiSquared[0])
		db = 1 + (alpha * e.oneOverRadiiSquared[1])
		dc = 1 + (alpha * e.oneOverRadiiSquared[2])

		da2 := da * da
		db2 := db * db
		dc2 := dc * dc

		da3 := da * da2
		db3 := db * db2
		dc3 := dc * dc2

		s = x2/(e.radiiSquared[0]*da2) +
			y2/(e.radiiSquared[1]*db2) +
			z2/(e.radiiSquared[2]*dc2) - 1

		dSdA = -2 *
			(x2/(e.radiiToTheFourth[0]*da3) +
				y2/(e.radiiToTheFourth[1]*db3) +
				z2/(e.radiiToTheFourth[2]*dc3))
	}

	return &vec3.T{pos[0] / da, pos[1] / db, pos[2] / dc}
}

// ToGeo3D takes an arbitrary point, and returns the latitude and longitude.
func (e *Ellipsoid) ToGeo3D(pos *vec3.T) *Geo3D {
	p := e.ScaleToGeodeticSurface(pos)
	h := vec3.Sub(pos, p)
	isPos := math.Signbit(vec3.Dot(&h, pos))
	height := h.Length()
	if isPos {
		height = -height
	}
	return &Geo3D{*e.ToGeo2D(p), height}
}
