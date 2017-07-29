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
	{"Twitter HQ", 37.777000, -122.416583, "9q8yym4fz"},
	{"Uber HQ", 37.775253, -122.417527, "9q8yykf2c"},
	{"Denver Airport", 39.855242, -104.672130, "9xjddpkjd"},
	{"JFK Airport", 40.641026, -73.777903, "dr5x1nkx4"},
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

// Valid bits range is 0 <= bits <= 64
func TestEncodeInt_badBits(t *testing.T) {
	assert.Panics(t, func() {
		EncodeInt(0.0, 0.0, -1)
	})

	assert.Panics(t, func() {
		EncodeInt(0.0, 0.0, 65)
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

	neighborhood := Neighborhood(0x3)
	assert.Equal(t, 9, len(neighborhood))
	for _, hash := range neighborhood {
		assert.True(t, expectedNeighborhood[hash])
	}
}

func TestNeighborhood_atEquator_evenHashLength(t *testing.T) {
	expectedNeighborhood := map[string]bool{
		"ebpbpc": true, "s00001": true, "s00003": true,
		"ebpbpb": true, "s00000": true, "s00002": true,
		"7zzzzz": true, "kpbpbp": true, "kpbpbr": true,
	}

	bits := 6 * 5 // 6-character base32 geohash; each base32 characters represents 5 bits
	h := EncodeInt(0.0, 0.0, bits)
	neighborhood := Neighborhood(h)

	assert.Equal(t, 9, len(neighborhood))

	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[toBase32(neighbor)])
	}

}

func TestNeighborhood_atEquator_oddHashLength(t *testing.T) {
	expectedNeighborhood := map[string]bool{
		"ebpbr": true, "s0002": true, "s0003": true,
		"ebpbp": true, "s0000": true, "s0001": true,
		"7zzzz": true, "kpbpb": true, "kpbpc": true,
	}

	bits := 5 * 5 // 5-character base32 geohash; each base32 characters represents 5 bits
	h := EncodeInt(0.0, 0.0, bits)
	neighborhood := Neighborhood(h)

	assert.Equal(t, 9, len(neighborhood))

	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[toBase32(neighbor)])
	}

}

func TestNeighborhood_9charGeoHash(t *testing.T) {
	expectedNeighborhood := map[string]bool{
		"9v6kpsezc": true, "9v6kpsezf": true, "9v6kpsezg": true,
		"9v6kpsez9": true, "9v6kpsezd": true, "9v6kpseze": true,
		"9v6kpsez3": true, "9v6kpsez6": true, "9v6kpsez7": true,
	}

	bits := 9 * 5 // 9-character base32 geohash; each base32 characters represents 5 bits
	h := EncodeInt(30.260415, -97.751107, bits)
	neighborhood := Neighborhood(h)

	assert.Equal(t, 9, len(neighborhood))

	for _, neighbor := range neighborhood {
		assert.True(t, expectedNeighborhood[toBase32(neighbor)])
	}

}

func Benchmark_EncodeInt(b *testing.B) {
	locationCount := len(locations)
	i := 0
	var encodedGeohash uint

	f := func() {
		location := locations[i]
		encodedGeohash = EncodeInt(location.Lat, location.Lng, 36)
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
