package parser

import (
	"unsafe"
)

const (
	BOOL = iota
	CHAR
	INT
	UINT
	FLOAT
	STRING
	UNKNOWN
)

func getType(t interface{}) int {
	switch t.(type) {
	case bool:
		return BOOL
	case rune:
		return CHAR
	case int:
		return INT
	case uint:
		return UINT
	case float64:
		return FLOAT
	case string:
		return STRING
	default:
		return UNKNOWN
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func RotateLeftUint(x uint, n uint) uint {
	size := uint(unsafe.Sizeof(x) * 8)
	n %= size
	return (x << n) | (x >> (size - n))
}

func RotateRightUint(x uint, n uint) uint {
	size := uint(unsafe.Sizeof(x) * 8)
	n %= size
	return (x >> n) | (x << (size - n))
}

func RotateLeftInt(x int, n uint) int {
	size := uint(unsafe.Sizeof(x) * 8)
	n %= size
	return (x << n) | (x >> (size - n))
}

func RotateRightInt(x int, n uint) int {
	size := uint(unsafe.Sizeof(x) * 8)
	n %= size
	return (x >> n) | (x << (size - n))
}

func toUint(value interface{}) uint {
	switch v := value.(type) {
	case bool:
		if v {
			return 1
		}
		return 0
	case int:
		if v < 0 {
			return 0
		}
		return uint(v)
	case uint:
		return v
	case float64:
		if v < 0.0 || v != float64(uint(v)) {
			return 0
		}
		return uint(v)
	default:
		return 0
	}
}
