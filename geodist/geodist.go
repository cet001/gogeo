package geodist

import (
	"github.com/cet001/mathext"
	"math"
)

// A geographical coordinate that represents a specific point on the earth.
type Coord struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

// Calculates the great-circle distance (in kilomerers) between geographical
// coordinates c1 and c2.
func Haversine(c1, c2 Coord) float64 {
	aLat, bLat := float64(c1.Lat)*mathext.Deg2rad, float64(c2.Lat)*mathext.Deg2rad
	deltaLat := bLat - aLat
	deltaLng := float64(c2.Lng-c1.Lng) * mathext.Deg2rad

	calc1 := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(aLat)*math.Cos(bLat)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	calc2 := 2 * math.Atan2(math.Sqrt(calc1), math.Sqrt(1-calc1))

	const earthRadiusKm float64 = 6371
	return earthRadiusKm * calc2
}

// Returns the Euclidean distance (in kilometers) between coordinates c1 and
// c2.  This distance function is much faster than either the Haversine formula
// or the spherical law of cosines, as it does not fully take into account the
// curvature of the earth (instead, the length of the longitude degree is
// calcluated based on the latitude).  As a result of this, its accuracy
// decreases as the distance between c1 and c2 increases.
//
// This function is practical in situations where you need to quickly calculate
// the distance between two "relatively close" points on a map (e.g. 2 houses
// located in the same neighborhood or city).
//
// See http://jonisalonen.com/2014/computing-distance-between-coordinates-can-be-simple-and-fast/
func Euclidean(c1, c2 Coord) float64 {
	const degreeLengthKm float64 = 110.25 // length of 1 degree longitude at the equator
	x := float64(c1.Lat - c2.Lat)
	y := float64(c1.Lng-c2.Lng) * math.Cos(float64(c2.Lat)*mathext.Deg2rad)
	return degreeLengthKm * math.Sqrt((x*x)+(y*y))
}
