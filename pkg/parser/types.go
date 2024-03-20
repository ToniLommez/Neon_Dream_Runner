package parser

const (
	BOOL = iota
	UINT
	INT
	FLOAT
	STRING
	UNKNOWN
)

func getType(t interface{}) int {
	switch t.(type) {
	case int:
		return INT
	case uint:
		return UINT
	case bool:
		return BOOL
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

/* func toInt(value any) (r int, err error) {
	switch v := value.(type) {
	case int:
		r = v
	case float64:
		r = int(v)
	case bool:
		if v {
			r = 1
		} else {
			r = 0
		}
	default:
		err = errutils.Error(-1, -1, fmt.Sprintf("%v", value), errutils.RUNTIME, fmt.Sprintf("Error converting %v to Int", value))
		r = 0
	}
	return
}

func toFloat(value any) (r float64, err error) {
	switch v := value.(type) {
	case int:
		r = float64(v)
	case float64:
		r = v
	case bool:
		if v {
			r = 1.0
		} else {
			r = 0.0
		}
	default:
		err = errutils.Error(-1, -1, fmt.Sprintf("%v", value), errutils.RUNTIME, fmt.Sprintf("Error converting %v to Float", value))
		r = 0
	}
	return
} */
