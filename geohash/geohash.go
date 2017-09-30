package geohash

import (
	"fmt"
	"unsafe"
)

const base32symbols string = "0123456789bcdefghjkmnpqrstuvwxyz"

// Encodes a (lat,lng) geo point as a N-bit geohash.
func Encode(lat, lng float32, bits int) uint {
	bits = validateGeoHashBitWidth(bits)

	// Adapted from https://www.factual.com/blog/how-geohashes-work
	var minLat, maxLat float64 = -90.0, 90.0
	var minLng, maxLng float64 = -180.0, 180.0

	var result uint = 0
	var lat64, lng64 = float64(lat), float64(lng)

	for i := 0; i < bits; i++ {
		if (i & 0x1) == 0 { // even bit: bisect longitude
			midpoint := (minLng + maxLng) / 2
			if lng64 < midpoint {
				result <<= 1      // push a zero bit
				maxLng = midpoint // shrink range downwards
			} else {
				result = result<<1 | 1 // push a one bit
				minLng = midpoint      // shrink range upwards
			}
		} else { // odd bit: bisect latitude
			midpoint := (minLat + maxLat) / 2
			if lat64 < midpoint {
				result <<= 1      // push a zero bit
				maxLat = midpoint // shrink range downwards
			} else {
				result = result<<1 | 1 // push a one bit
				minLat = midpoint      // shrink range upwards
			}
		}
	}
	return result
}

// Encodes a (lat,lng) geo point as a base-32 geohash string of the specified
// length (not to exceed 12 characters).
func EncodeBase32(lat, lng float32, length int) string {
	length = validateBase32GeoHashLength(length)
	hashBits := length * 5
	return toBase32(Encode(lat, lng, hashBits))
}

// Returns a slice containing the provided geohash along with its 8 surrounding
// geohash tiles.
func Neighborhood(lat, lng float32, bits int) []uint {
	h := Encode(lat, lng, validateGeoHashBitWidth(bits))

	// adapted from: https://github.com/yinqiwen/geohash-int/blob/b01291be60015cd399227f2e3305c5a3262f68c1/geohash.c
	return []uint{
		h,                       // me
		moveY(h, 1),             // north neighbor
		moveY(h, -1),            // south neighbor
		moveX(h, 1),             // east neighbor
		moveX(h, -1),            // west neighbor
		moveX(moveY(h, 1), 1),   // northeast neighbor
		moveX(moveY(h, -1), 1),  // southeast neighbor
		moveX(moveY(h, 1), -1),  // northwest neighbor
		moveX(moveY(h, -1), -1), // southwest neighbor
	}
}

func NeighborhoodBase32(lat, lng float32, length int) []string {
	length = validateBase32GeoHashLength(length)
	bits := length * 5
	neighborhood := Neighborhood(lat, lng, bits)
	neighborhoodBase32 := make([]string, len(neighborhood))
	for i, neighbor := range neighborhood {
		neighborhoodBase32[i] = toBase32(neighbor)
	}

	return neighborhoodBase32
}

// Neighborhood() helper that calculates the east-west adjacent tile based on
// the 'dir' arg (west = -1, east = 1).
func moveX(geohash uint, dir int) uint {
	var x uint = geohash & 0xAAAAAAAAAAAAAAAA
	var y uint = geohash & 0x5555555555555555
	const zz uint = 0x5555555555555555

	if dir > 0 {
		x += (zz + 1)
	} else {
		x |= zz
		x -= (zz + 1)
	}

	x &= 0xAAAAAAAAAAAAAAAA
	return x | y
}

// Neighborhood() helper that calculates the north-south adjacent tile based on
// the 'dir' arg (north = 1, south = -1).
func moveY(geohash uint, dir int) uint {
	var x uint = geohash & 0xAAAAAAAAAAAAAAAA
	var y uint = geohash & 0x5555555555555555
	const zz uint = 0xAAAAAAAAAAAAAAAA

	if dir > 0 {
		y += (zz + 1)
	} else {
		y |= zz
		y -= (zz + 1)
	}

	y &= 0x5555555555555555
	return x | y
}

// Convert the N-bit geohash integer into a base32 string.
func toBase32(h uint) string {
	// Pre-calculate how many base-32 characters we'll need
	hTmp := h
	charCount := 0
	for hTmp > 0 {
		hTmp >>= 5
		charCount++
	}

	// generate the geohash string
	b := make([]byte, charCount)
	for i := (charCount - 1); i >= 0; i-- {
		b[i] = base32symbols[h&0x1F]
		h >>= 5
	}

	// Performance optimization: eliminate []byte copy when converting to string
	return *(*string)(unsafe.Pointer(&b))
}

func validateGeoHashBitWidth(bits int) int {
	if bits < 0 || bits > 64 {
		panic(fmt.Sprintf("'bits' must be in the range [0, 64]. Actual value was %v", bits))
	}
	return bits
}

func validateBase32GeoHashLength(length int) int {
	if length < 1 || length > 12 {
		panic(fmt.Errorf("Requested geohash length out of range: %v. Valid range is [1..12]", length))
	}
	return length
}
