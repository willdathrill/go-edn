package edn

import (
	//"hash/adler32" depends on testing for hash performance
	hash "hash/fnv1" // named to be able to switch easily
	"math"
	r "reflect"
)

// this function is only for scalar values:
// not arrays, slices, pointers, chans, maps, or structs
// however, string, and []byte are accepted
func hashScalar(in interface{}) int32 {
	if in == nil { return 0 }
	var to []byte
	switch i :=in {
	case byte, uint8, int8:
		to = []byte{byte(i)}
	case uint16, int16:
		u := uint16(i)
		to = []byte{byte(u>>8), byte(u)}
	case uint32, int32, uint, int, rune:
		u := uint32(i)
		to = []byte{byte(u>>24),byte(u>>16),byte(u>>8),byte(u)}
	case uint64, int64:
		u := uint64(i)
		to = []byte{
			byte(u>>56),byte(u>>48),byte(u>>40),byte(u>>32),
			byte(u>>24),byte(u>>16),byte(u >> 8),byte(u)
		}
	case float32:
		to = toBytes(math.Float32bits(i))
	case float64:
		to = toBytes(math.Float64bits(i))
	case complex64, complex128:
		re, im := real(i), imag(i)
		to = append(toBytes(re),toBytes(im)...)
	case []byte:
		to = i
	case string:
		to = []byte(i)
	}
	s := hash.New32()
	s.Write(to)
	return s.Sum32()
}
	