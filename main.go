package gogeo

import (
	"fmt"
	_ "github.com/cet001/gogeo/geohash"
	"strconv"
	"strings"
)

// Parses a latitude/longitude point of the form "45.521066,-122.684984" into a
// pair of float32 values.
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
		return lat, lng, fmt.Errorf("Latitude not in [-90, 90] range: %v", lat)
	}

	if lng < -180.0 || lng > 180.0 {
		return lat, lng, fmt.Errorf("Longitude not in [-180, 180] range: %v", lng)
	}

	return lat, lng, nil
}

func parseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	return float32(f), err
}
