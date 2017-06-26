package gogeo

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLatLng(t *testing.T) {
	lat, lng, err := ParseLatLng("50.456, -80.2468")
	assert.Nil(t, err)
	assert.Equal(t, lat, float32(50.456))
	assert.Equal(t, lng, float32(-80.2468))
}

func TestParseLatLng_badInputs(t *testing.T) {
	badInputs := []string{
		"",
		"   ",
		" ,  ",
		"aaa,12.345",          // lat not a number
		"12.345,bbb",          // lng not a number
		"aaa,bbb",             // both lat and lng not numbers
		"11.11, 22.22, 33.33", // too many items
		"11.11",               // not enough items
		"-91.1, 0.0",          // latitude value too small. not in range [-90, +90]
		"90.1, 0.0",           // latitude value too large. not in range [-90, +90]
		"0.0, -181.1",         // longitude value too small. not in range [-180, +180]
		"0.0, +181.1",         // longitude value too large. not in range [-180, +180]
	}

	for _, badInput := range badInputs {
		_, _, err := ParseLatLng(badInput)
		assert.NotNil(t, err)
	}
}
