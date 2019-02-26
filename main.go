package gogeo

import (
	"fmt"
	_ "github.com/cet001/gogeo/geohash"
	"strconv"
	"strings"
)

// Parses a string representing a latitude value in degrees into a float32 value,
// returning an error if the input was either not parseable as a float32 or was
// not in the valid latitude range of [-90 .. +90].
func ParseLat(latitude string) (float32, error) {
	lat, err := parseFloat32(latitude)
	if err != nil {
		return 0, fmt.Errorf(`Error parsing latitude as float32: "%v"`, err)
	}
	if lat < -90.0 || lat > 90.0 {
		return lat, fmt.Errorf("%v is not a valid latitude (valid range is[-90, 90])", lat)
	}
	return lat, nil
}

// Parses a string representing a longitude value in degrees into a float32 value,
// returning an error if the input was either not parseable as a float32 or was
// not in the valid longitude range of [-180 .. +180].
func ParseLng(longitude string) (float32, error) {
	lng, err := parseFloat32(longitude)
	if err != nil {
		return 0, fmt.Errorf(`Error parsing longitude as float32: "%v"`, err)
	}
	if lng < -180.0 || lng > 180.0 {
		return lng, fmt.Errorf("%v is not a valid longitude (valid range is[-180, 180])", lng)
	}
	return lng, nil
}

// Parses a formatted latitude/longitude coordinate of the form "<lat>,<lng>"
// into a pair of float32 values representing the latitude and longitude.
func ParseLatLng(latLng string) (float32, float32, error) {
	latLngSlice := strings.Split(latLng, ",")
	if len(latLngSlice) != 2 {
		return 0, 0, fmt.Errorf(`Can't parse LatLng string: "%v"`, latLng)
	}

	lat, err := ParseLat(latLngSlice[0])
	if err != nil {
		return 0, 0, err
	}

	lng, err := ParseLng(latLngSlice[1])
	if err != nil {
		return 0, 0, err
	}

	return lat, lng, nil
}

func parseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	return float32(f), err
}
