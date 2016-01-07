package globe

import "math"

// Geo2D holds the latitude and longitude of a coordinate on a globe's surface.
type Geo2D struct {
	lat, lon float64
	isRad    bool
}

// NewGeo2D takes in the latitude and longitude of a coordinate in degrees, and creates a Geo2D type to represent it.
func NewGeo2D(lat, lon float64) *Geo2D {
	return &Geo2D{lat, lon, false}
}

// ToRad creates a new Geo2D with the coordinate of the original Geo2D converted to radians. It will simply copy the data if it is already in radians.
func (g *Geo2D) ToRad() *Geo2D {
	newGeo := *g
	if !newGeo.isRad {
		newGeo.lat = newGeo.lat * math.Pi / 180
		newGeo.lon = newGeo.lon * math.Pi / 180
		newGeo.isRad = true
	}
	return &newGeo
}

// ToDeg creates a new Geo2D with the coordinate of the original Geo2D converted to degrees. It will simply copy the data if it is already in degrees.
func (g *Geo2D) ToDeg() *Geo2D {
	newGeo := *g
	if newGeo.isRad {
		newGeo.lat = newGeo.lat * 180 / math.Pi
		newGeo.lon = newGeo.lon * 180 / math.Pi
		newGeo.isRad = false
	}
	return &newGeo
}

// LatLon returns the latitude and longitude of the coordinate represented by the Geo2D.
func (g *Geo2D) LatLon() (float64, float64) {
	return g.lat, g.lon
}

// Geo3D holds the latitude and longitude and height of a coordinate relative to a globe's surface.
type Geo3D struct {
	Geo2D
	height float64
}

// NewGeo3D takes in the latitude and longitude of a coordinate in degrees, and the height from the surface, and creates a Geo3D type to represent it.
func NewGeo3D(lat, lon, height float64) *Geo3D {
	return &Geo3D{*NewGeo2D(lat, lon), height}
}

// ToRad creates a new Geo3D with the coordinate of the original Geo3D converted to radians. It will simply copy the data if it is already in radians.
func (g *Geo3D) ToRad() *Geo3D {
	newGeo := *g
	newGeo.Geo2D = *newGeo.Geo2D.ToRad()
	return &newGeo
}

// ToDeg creates a new Geo3D with the coordinate of the original Geo3D converted to degrees. It will simply copy the data if it is already in degrees.
func (g *Geo3D) ToDeg() *Geo3D {
	newGeo := *g
	newGeo.Geo2D = *newGeo.Geo2D.ToDeg()
	return &newGeo
}

// GetHeight returns the height of the coordinate represented by the Geo3D.
func (g *Geo3D) GetHeight() float64 {
	return g.height
}
