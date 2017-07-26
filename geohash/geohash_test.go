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

func TestEncode(t *testing.T) {
	for _, location := range locations {
		for i := 1; i <= 9; i++ {
			h, err := Encode(location.Lat, location.Lng, i)
			assert.Nil(t, err)
			assert.Equal(t, location.Geohash[:i], h, fmt.Sprintf("%v [length=%v]", location.Name, i))
		}
	}
}

func TestEncode_badLength(t *testing.T) {
	_, err := Encode(0.0, 0.0, 0) // geohash length too small
	assert.NotNil(t, err)

	_, err = Encode(0.0, 0.0, 13) // geohash length too big
	assert.NotNil(t, err)
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

var encodedGeohash string

func Benchmark_Encode(b *testing.B) {
	locationCount := len(locations)
	i := 0

	f := func() {
		location := locations[i]
		encodedGeohash, _ = Encode(location.Lat, location.Lng, 7)
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
