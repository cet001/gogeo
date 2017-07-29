package geohash

import (
	"fmt"
	"unsafe"
)

const base32symbols string = "0123456789bcdefghjkmnpqrstuvwxyz"

// Encodes a (lat,lng) geo point as a N-bit integer.
func EncodeInt(lat, lng float32, bits int) uint {
	if bits < 0 || bits > 64 {
		panic(fmt.Sprintf("'bits' must be in the range [0, 64]. Actual value was %v", bits))
	}

	// Adapted from https://www.factual.com/blog/how-geohashes-work
	var minLat, maxLat float32 = -90.0, 90.0
	var minLng, maxLng float32 = -180.0, 180.0

	var result uint = 0

	for i := 0; i < bits; i++ {
		if (i & 0x1) == 0 { // even bit: bisect longitude
			midpoint := (minLng + maxLng) / 2
			if lng < midpoint {
				result <<= 1      // push a zero bit
				maxLng = midpoint // shrink range downwards
			} else {
				result = result<<1 | 1 // push a one bit
				minLng = midpoint      // shrink range upwards
			}
		} else { // odd bit: bisect latitude
			midpoint := (minLat + maxLat) / 2
			if lat < midpoint {
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
	if length < 1 || length > 12 {
		panic(fmt.Errorf("Requested geohash length out of range: %v. Valid range is [1..12]", length))
	}

	hashBits := length * 5
	return toBase32(EncodeInt(lat, lng, hashBits))
}

// Returns a slice containing the provided geohash along with its 8 surrounding
// geohash tiles.
func Neighborhood(geohash uint, bits uint) []uint {
	// round bits up to nearest even integer
	if bits&0x1 == 1 {
		bits += 1
	}

	return []uint{
		geohash,                                   // me
		moveY(geohash, bits, 1),                   // north neighbor
		moveY(geohash, bits, -1),                  // south neighbor
		moveX(geohash, bits, 1),                   // east neighbor
		moveX(geohash, bits, -1),                  // west neighbor
		moveX(moveY(geohash, bits, 1), bits, 1),   // northeast neighbor
		moveX(moveY(geohash, bits, -1), bits, 1),  // southeast neighbor
		moveX(moveY(geohash, bits, 1), bits, -1),  // northwest neighbor
		moveX(moveY(geohash, bits, -1), bits, -1), // southwest neighbor
	}
}

// Convert the N-bit geohash value into a base32 string.
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

// Neighborhood() helper that calculates the east-west adjacent tile based on
// the 'dir' arg (west = -1, east = 1).
func moveX(geohash, bits uint, dir int) uint {

	if dir == 0 {
		return geohash
	}

	var x uint = geohash & 0xAAAAAAAAAAAAAAAA
	var y uint = geohash & 0x5555555555555555
	var zz uint = 0x5555555555555555 >> (64 - bits)

	if dir > 0 {
		x += (zz + 1)
	} else {
		x |= zz
		x -= (zz + 1)
	}

	x &= (0xAAAAAAAAAAAAAAAA >> (64 - bits))
	return x | y
}

// Neighborhood() helper that calculates the north-south adjacent tile based on
// the 'dir' arg (north = 1, south = -1).
func moveY(geohash, bits uint, dir int) uint {
	if dir == 0 {
		return geohash
	}

	var x uint = geohash & 0xAAAAAAAAAAAAAAAA
	var y uint = geohash & 0x5555555555555555
	var zz uint = 0xAAAAAAAAAAAAAAAA >> (64 - bits)

	if dir > 0 {
		y += (zz + 1)
	} else {
		y |= zz
		y -= (zz + 1)
	}

	y &= (0x5555555555555555 >> (64 - bits))
	return x | y
}
