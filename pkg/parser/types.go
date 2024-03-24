package parser

import (
	"unsafe"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

const (
	UNKNOWN = iota
	BOOL
	CHAR
	INT
	UINT
	FLOAT
	STRING
	NIL
	UNDEFINED
)

func getType(t any) int {
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
	case nil:
		return NIL
	default:
		return UNKNOWN
	}
}

func tokenToType(t lexer.Token) int {
	switch t.Type {
	case lexer.BOOL:
		return BOOL
	case lexer.CHAR:
		return CHAR
	case lexer.INT:
		return INT
	case lexer.UINT:
		return UINT
	case lexer.FLOAT:
		return FLOAT
	case lexer.STRING:
		return STRING
	case lexer.NIL:
		return NIL
	case lexer.UNDEFINED:
		return UNDEFINED
	default:
		return UNKNOWN
	}
}

func typeToString(t int) string {
	switch t {
	case UNKNOWN:
		return "UNKNOWN"
	case BOOL:
		return "BOOL"
	case CHAR:
		return "CHAR"
	case INT:
		return "INT"
	case UINT:
		return "UINT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case NIL:
		return "NIL"
	case UNDEFINED:
		return "UNDEFINED"
	default:
		return "this should never be printed, errorcode: 3286"
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

func toUint(value any) uint {
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

/* func cast(value any, targetType int) (any, bool) {
	switch v := value.(type) {
	case bool:
		switch targetType {
		case BOOL:
			return v, true
		case CHAR:
			return v, false
		case INT:
			return utils.Ternary(v, int(1), int(0)), true
		case UINT:
			return utils.Ternary(v, uint(1), uint(0)), true
		case FLOAT:
			return utils.Ternary(v, float64(1), float64(0)), true
		case STRING:
			return v, false
		default:
			return v, false
		}
	case rune:
		switch targetType {
		case BOOL:
			return v != 0, true
		case CHAR:
			return v, true
		case INT:
			return int(v), true
		case UINT:
			return uint(v), true
		case FLOAT:
			return float64(v), true
		case STRING:
			return v, false
		default:
			return v, false
		}
	case int:
		switch targetType {
		case BOOL:
			return v != 0, true
		case CHAR:
			return rune(v), true
		case INT:
			return v, true
		case UINT:
			return v, false
		case FLOAT:
			return float64(v), true
		case STRING:
			return v, false
		default:
			return v, false
		}
	case uint:
		switch targetType {
		case BOOL:
			return v != 0, true
		case CHAR:
			return rune(v), true
		case INT:
			return v, false
		case UINT:
			return v, true
		case FLOAT:
			return float64(v), true
		case STRING:
			return v, false
		default:
			return v, false
		}
	case float64:
		switch targetType {
		case BOOL:
			return v, false
		case CHAR:
			return v, false
		case INT:
			return v, false
		case UINT:
			return v, false
		case FLOAT:
			return v, true
		case STRING:
			return v, false
		default:
			return v, false
		}
	case string:
		switch targetType {
		case BOOL:
			return 0, false
		case CHAR:
			if len(v) == 1 {
				return rune(v[0]), true
			} else {
				return v, false
			}
		case INT:
			return v, false
		case UINT:
			return v, false
		case FLOAT:
			return v, false
		case STRING:
			return v, true
		default:
			return v, false
		}
	default:
		return v, false
	}
} */
