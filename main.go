package gogeo

import (
	"fmt"
	_ "github.com/cet001/gogeo/geohash"
	"strconv"
	"strings"
)

// Parses a formatted latitude/longitude geo point of the form "<lat>,<lng>" into
// a pair of float32 values.
func ParseLatLng(latLng string) (float32, float32, error) {
	latLngSlice := strings.Split(latLng, ",")
	if len(latLngSlice) != 2 {
		return 0, 0, fmt.Errorf("Can't parse LatLng string: \"%v\"", latLng)
	}

	lat, err := parseFloat32(latLngSlice[0])
	if err != nil {
		return 0, 0, err
	}

	lng, err := parseFloat32(latLngSlice[1])
	if err != nil {
		return 0, 0, err
	}

	if lat < -90.0 || lat > 90.0 {
		return lat, lng, fmt.Errorf("%v is not a valid latitude (valid range is[-90, 90])", lat)
	}

	if lng < -180.0 || lng > 180.0 {
		return lat, lng, fmt.Errorf("%v is not a valid longitude (valid range is[-180, 180])", lng)
	}

	return lat, lng, nil
}

func parseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	return float32(f), err
}
