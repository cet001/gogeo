package geohash

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestLocation struct {
	Name    string
	Lat     float32
	Lng     float32
	Geohash string
}

// List of known geo points
var locations = []TestLocation{
	{"Twitter HQ", 37.777000, -122.416583, "9q8yym4fy"},
	{"Uber HQ", 37.775253, -122.417527, "9q8yykf2b"},
	{"Denver Airport", 39.855242, -104.672130, "9xjddpkjd"},
	{"JFK Airport", 40.641026, -73.777903, "dr5x1nkx4"},
}

func ExampleEncode() {
	fmt.Println(Encode(37.751223, -122.438297, 8))
	fmt.Println(Encode(37.751223, -122.438297, 32))
	// Output:
	// 77
	// 1301409192
}

// Valid bits range is 0 <= bits <= 64
func TestEncode_badBits(t *testing.T) {
	assert.Panics(t, func() {
		Encode(0.0, 0.0, -1)
	})

	assert.Panics(t, func() {
		Encode(0.0, 0.0, 65)
	})
}

func ExampleEncodeBase32() {
	fmt.Println(EncodeBase32(37.751223, -122.438297, 6))
	fmt.Println(EncodeBase32(37.751223, -122.438297, 5))
	fmt.Println(EncodeBase32(37.751223, -122.438297, 4))
	// Output:
	// 9q8yvb
	// 9q8yv
	// 9q8y
}

func TestEncodeBase32(t *testing.T) {
	for _, location := range locations {
		for i := 1; i <= 9; i++ {
			h := EncodeBase32(location.Lat, location.Lng, i)
			assert.Equal(t, location.Geohash[:i], h, fmt.Sprintf("%v [length=%v]", location.Name, i))
		}
	}
}

func TestEncodeBase32_badLength(t *testing.T) {
	assert.Panics(t, func() {
		EncodeBase32(0.0, 0.0, 0) // geohash length too small
	})

	assert.Panics(t, func() {
		EncodeBase32(0.0, 0.0, 13) // geohash length too big
	})
}

func TestNeighborhood(t *testing.T) {
	expectedNeighborhood := map[uint]bool{
		0x3: true, // me ('0011' 4-bit binary)
		0x2: true, // North neighbor
		0x6: true, // South neighbor
		0x9: true, // East neighbor
		0x1: true, // West neighbor
		0x8: true, // NE neighbor
		0xC: true, // SE neighbor
		0x0: true, // NW neighbor
		0x4: true, // SW neighbor
	}

	// Note: The geo point (0,-90) corresonds to the geohash 0x3
	neighborhood := Neighborhood(-10.0, -90.0, 4)

	assert.Equal(t, 9, len(neighborhood))
	for _, hash := range neighborhood {
		assert.True(t, expectedNeighborhood[hash])
	}
}

func TestNeighborhoodBase32_atEquator_evenHashLength(t *testing.T) {
	// These results generated with http://www.movable-type.co.uk/scripts/geohash.html
	expectedNeighborhood := map[string]bool{
		"ebpbpc": true, "s00001": true, "s00003": true,
		"ebpbpb": true, "s00000": true, "s00002": true,
		"7zzzzz": true, "kpbpbp": true, "kpbpbr": true,
	}

	neighborhood := NeighborhoodBase32(0.0, 0.0, 6)
	assert.Equal(t, 9, len(neighborhood))
	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[neighbor])
	}
}

func TestNeighborhoodBase32_atEquator_oddHashLength(t *testing.T) {
	// These results generated with http://www.movable-type.co.uk/scripts/geohash.html
	expectedNeighborhood := map[string]bool{
		"ebpbr": true, "s0002": true, "s0003": true,
		"ebpbp": true, "s0000": true, "s0001": true,
		"7zzzz": true, "kpbpb": true, "kpbpc": true,
	}

	neighborhood := NeighborhoodBase32(0.0, 0.0, 5)
	assert.Equal(t, 9, len(neighborhood))
	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[neighbor])
	}

}

func TestNeighborhoodBase32_9charGeoHash(t *testing.T) {
	// These results generated with http://www.movable-type.co.uk/scripts/geohash.html
	expectedNeighborhood := map[string]bool{
		"9v6kpsezc": true, "9v6kpsezf": true, "9v6kpsezg": true,
		"9v6kpsez9": true, "9v6kpsezd": true, "9v6kpseze": true,
		"9v6kpsez3": true, "9v6kpsez6": true, "9v6kpsez7": true,
	}

	neighborhood := NeighborhoodBase32(30.260415, -97.751107, 9)
	assert.Equal(t, 9, len(neighborhood))
	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[neighbor])
	}

}

func Benchmark_Encode(b *testing.B) {
	locationCount := len(locations)
	i := 0
	var encodedGeohash uint

	f := func() {
		location := locations[i]
		encodedGeohash = Encode(location.Lat, location.Lng, 36)
		i++
		if i == locationCount {
			i = 0
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f()
	}
}

func Benchmark_EncodeBase32(b *testing.B) {
	locationCount := len(locations)
	i := 0
	var encodedGeohash string

	f := func() {
		location := locations[i]
		encodedGeohash = EncodeBase32(location.Lat, location.Lng, 7)
		i++
		if i == locationCount {
			i = 0
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f()
	}
}
