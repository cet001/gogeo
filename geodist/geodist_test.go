package geodist

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Segment struct {
	a    Coord
	b    Coord
	dist float64
}

var taraIndianCuisine = Coord{37.765492, -122.431229}
var blackbirdBar = Coord{37.767487, -122.429633}
var madDogInTheFog = Coord{37.772528, -122.431030}
var spinCityCoffee = Coord{37.749287, -122.429181}
var civicCenterBart = Coord{37.780396, -122.414135}
var monteryBayAquarium = Coord{36.624503, -121.901581}
var modesto = Coord{37.681253, -120.997561}
var sfo = Coord{37.621313, -122.378945}
var lax = Coord{33.941446, -118.408594}

var testSegments = []Segment{
	Segment{taraIndianCuisine, blackbirdBar, 0.24},
	Segment{taraIndianCuisine, madDogInTheFog, 0.76},
	Segment{taraIndianCuisine, spinCityCoffee, 1.83},
	Segment{taraIndianCuisine, civicCenterBart, 2.17},
	Segment{sfo, monteryBayAquarium, 120.0},
	Segment{sfo, modesto, 122.0},
	Segment{sfo, lax, 543.0},
}

func ExampleApproxDist() {
	voodooDoughnuts := Coord{Lat: 45.522869, Lng: -122.673132}
	powellsCityOfBooks := Coord{Lat: 45.523437, Lng: -122.681381}
	d := ApproxDist(voodooDoughnuts, powellsCityOfBooks)
	fmt.Printf("Voodoo Doughnuts is %.2f km from Powell's City of Books.\n", d)
	// Output:
	// Voodoo Doughnuts is 0.64 km from Powell's City of Books.
}

func TestApproxDist(t *testing.T) {
	for _, testSegment := range testSegments {
		distKm := ApproxDist(testSegment.a, testSegment.b)
		assert.InDelta(t, testSegment.dist, distKm, deltaErr(testSegment.dist))
	}
}

func TestHaversine(t *testing.T) {
	for _, testSegment := range testSegments {
		distKm := Haversine(testSegment.a, testSegment.b)
		assert.InDelta(t, testSegment.dist, distKm, deltaErr(testSegment.dist))
	}
}

// Based on the expected distance, return the acceptable delta error of the
// expected vs. actual distance (in Km).
func deltaErr(distKm float64) float64 {
	switch {
	case distKm < 10:
		return 0.1
	case distKm < 100:
		return distKm * 0.01
	default:
		return distKm * 0.02
	}
}
