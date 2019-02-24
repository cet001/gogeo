package gogeo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleParseLatLng() {
	lat, lng, _ := ParseLatLng("50.456, -80.2468")
	fmt.Println(lat)
	fmt.Println(lng)
	// Output:
	// 50.456
	// -80.2468
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

func TestParseFloat32(t *testing.T) {
	testCases := []struct {
		s     string
		f     float32
		isErr bool
	}{
		{s: "0", f: 0.0},
		{s: "1.2345", f: 1.2345},
		{s: "-1.2345", f: -1.2345},
		{s: "  1.5", f: 1.5},
		{s: "ABC", isErr: true},
		{s: "1x", isErr: true},
		{s: "", isErr: true},
	}

	for i, testCase := range testCases {
		testCaseLabel := fmt.Sprintf("testCases[%v]", i)
		f, err := parseFloat32(testCase.s)
		if testCase.isErr {
			assert.NotNil(t, err, "%v should have returned an error", testCaseLabel)
		} else {
			assert.Nil(t, err, testCaseLabel)
			assert.Equal(t, testCase.f, f, testCaseLabel)
		}
	}
}
