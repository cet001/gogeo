package geohash

import (
	"fmt"
	"unsafe"
)

const base32 string = "0123456789bcdefghjkmnpqrstuvwxyz"

// Encodes a (lat,lng) geo point as a base-32 geohash string of the specified length
// (not to exceed 12 characters).
func Encode(lat, lng float32, length int) (string, error) {
	if length < 1 || length > 12 {
		return "", fmt.Errorf("Requested geohash length out of range: %v. Valid range is [1..12]", length)
	}

	hashBits := length * 5
	h := EncodeInt(lat, lng, hashBits)

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
		b[i] = base32[h&0x1F]
		h >>= 5
	}

	// Performance optimization: eliminate []byte copy when converting to string
	return *(*string)(unsafe.Pointer(&b)), nil
}

// Encodes a (lat,lng) geo point as a N-bit integer.
func EncodeInt(lat, lng float32, bits int) uint {
	if bits < 0 || bits > 64 {
		panic(fmt.Sprintf("bits arg mus be in the range [0, 64]. Actual value was %v", bits))
	}

	// Adapted from https://www.factual.com/blog/how-geohashes-work
	var minLat, maxLat float32 = -90.0, 90.0
	var minLng, maxLng float32 = -180.0, 180.0

	var result uint = 0

	for i := 0; i < bits; i++ {
		if i%2 == 0 { // even bit: bisect longitude
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
